package executor

import (
	"fmt"
	"time"

	"github.com/skyflow-workflow/skyflow_backbend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

// Choice Choice状态
type Choice struct {
	*ExecutionStep
	State *states.ChoiceState
}

// NewChoiceFromID NewChoiceFromID
func NewChoiceFromID(id int, executor *Executor) (*Choice, error) {

	var err error

	var dbStep = &po.Step{}

	dbStep, err = executor.ExecutionService.QueryStepByID(id, []string{}, nil)
	if err != nil {
		return nil, err
	}

	state, err := NewChoiceFromData(dbStep, executor)

	return state, err
}

// NewChoiceFromData NewChoiceFromData
func NewChoiceFromData(dbStep *po.Step, executor *Executor) (*Choice, error) {

	var err error

	baseStep, err := NewExecutionStep(dbStep, executor)
	if err != nil {
		return nil, err
	}
	state, ok := baseStep.State.(*states.ChoiceState)
	if !ok {
		return nil, fmt.Errorf("step state is not 'Pass' type")
	}
	choiceState := &Choice{
		ExecutionStep: baseStep,
		State:         state,
	}
	return choiceState, err

}

// GetNextState  Get Next State
func (c *Choice) GetNextState() (NextStep, error) {
	var sns NextStep
	var err error
	var input interface{}
	input, err = toolkit.DecodeStringToMap(c.Data.Input)
	if err != nil {
		return sns, err
	}
	ns, err := c.State.GetNextState(input)
	if err != nil {
		return sns, err
	}
	sns = NextStep{
		NextState: ns,
		GroupID:   c.Data.GroupID,
	}
	return sns, err
}

// Run  Run
func (c *Choice) Run(_ queue.InnerMessageBody) error {
	var err error
	dbstate := c.Data

	starttime := time.Now()

	sns, err := c.GetNextState()
	if err != nil {
		return err
	}
	var outputstr string
	outputstr, err = toolkit.ToString(sns.Output)
	if err != nil {
		return err
	}

	tx, maker := c.ExecutionService.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)
	var updatestep po.Step

	now := time.Now()

	updatestep = po.Step{
		Status:     string(StepStatus.Success),
		Output:     outputstr,
		FinishTime: &now,
	}

	err = tx.Where(po.Step{ID: dbstate.ID}).Updates(&updatestep).Error
	if err != nil {

		return err
	}
	tx.Commit()

	event1 := vo.ExecutionEvent{
		ExecutionID: dbstate.ExecutionID,
		StartTime:   starttime,
		FinishTime:  now,
		StepID:      dbstate.ID,
		StepName:    dbstate.Name,
		Data: EventContent_ChoiceStateExecuted{
			Next:   sns.Name,
			Output: sns.Output,
		},
	}
	c.ExecutionService.SendExecutionEvents(event1)

	fns := FindNextStep{
		Name:    sns.Name,
		GroupID: sns.GroupID,
	}
	// 写入message queue
	message := NewStepMessage(dbstate.ExecutionID, MessageType.FindNextStep, dbstate.ID, fns)
	err = c.ExecutionService.SendInnerMessage(message, nil)
	if err != nil {
		return err
	}

	return nil
}
