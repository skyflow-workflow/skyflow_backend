package executor

import (
	"fmt"
	"time"

	"github.com/skyflow-workflow/skyflow_backend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backend/workflow/vo"
)

// Pass Pass ExecutionState
type Pass struct {
	*ExecutionStep // 继承 ExecutionStep
	State          *states.PassState
}

// NewPassFromID NewPassFromID
func NewPassFromID(id int, executor *Executor) (*Pass, error) {

	var err error

	dbStep, err := executor.ExecutionService.QueryStepByID(id, []string{}, nil)
	if err != nil {
		return nil, err
	}

	state, err := NewPassFromData(dbStep, executor)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func NewPassFromData(dbStep *po.Step, executor *Executor) (*Pass, error) {

	baseStep, err := NewExecutionStep(dbStep, executor)
	if err != nil {
		return nil, err
	}
	state, ok := baseStep.State.(*states.PassState)
	if !ok {
		return nil, fmt.Errorf("step state is not 'Pass' type")
	}
	pass := &Pass{
		ExecutionStep: baseStep,
		State:         state,
	}
	return pass, nil
}

// IsEnd IsEnd
func (p *Pass) IsEnd() bool {
	return p.State.IsEnd()
}

// GetResult GetResult
func (p *Pass) GetResult() (interface{}, error) {
	input, err := p.DecodeInput()
	if err != nil {
		return nil, err
	}
	res, err := p.State.GetResult(input)
	return res, err
}

// Run Run
func (p *Pass) Run(_ queue.InnerMessageBody) error {
	var err error
	var dbstate = p.Data

	starttime := time.Now()

	result, err := p.GetResult()
	if err != nil {
		return err
	}
	sns, err := p.ExecutionStep.GetNextStep(result)
	if err != nil {
		return err
	}
	outputStr, err := toolkit.ToString(sns.Output)
	if err != nil {
		return err
	}

	now := time.Now()

	tx, maker := p.ExecutionService.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	updatetask := po.Step{
		Output:     outputStr,
		FinishTime: &now,
		Status:     string(StepStatus.Success),
	}

	err = tx.Where(po.Step{ID: dbstate.ID}).Updates(&updatetask).Error
	if err != nil {
		return err
	}
	tx.Commit()

	// send event
	event1 := vo.ExecutionEvent{
		ExecutionID: dbstate.ExecutionID,
		StartTime:   starttime,
		FinishTime:  now,
		StepID:      dbstate.ID,
		StepName:    dbstate.Name,
		Data: EventContent_PassStateExecuted{
			Result: result,
		},
	}

	event2 := vo.ExecutionEvent{
		ExecutionID: dbstate.ExecutionID,
		StartTime:   starttime,
		FinishTime:  now,
		StepID:      dbstate.ID,
		StepName:    dbstate.Name,
		Data: EventContent_StateExited{
			Output: sns.Output,
		},
	}
	p.ExecutionService.SendExecutionEvents(event1, event2)

	fns := FindNextStep{
		Name:    sns.Name,
		GroupID: sns.GroupID,
	}
	// message queue send create message
	message := NewStepMessage(dbstate.ExecutionID, MessageType.FindNextStep, dbstate.ID, fns)
	err = p.ExecutionService.InnerQueue.SendInnerMessage(message, nil)
	if err != nil {
		return err
	}

	return nil
}
