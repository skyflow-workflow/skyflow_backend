package executor

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backend/workflow/vo"
)

// Suspend  executor for suspend state
type Suspend struct {
	*ExecutionStep
	State *states.SuspendState
}

// NewWaitFromID NewWaitFromID
func NewSuspendFromID(id int, executor *Executor) (*Suspend, error) {

	var err error

	var dbStep = &po.Step{}

	dbStep, err = executor.ExecutionService.QueryStepByID(id, []string{}, nil)
	if err != nil {
		return nil, err
	}

	state, err := NewSuspendFromData(dbStep, executor)

	return state, err

}

// NewSuspendFromData NewSuspendFromData
func NewSuspendFromData(dbStep *po.Step, executor *Executor) (*Suspend, error) {
	baseStep, err := NewExecutionStep(dbStep, executor)
	if err != nil {
		return nil, err
	}
	state, ok := baseStep.State.(*states.SuspendState)
	if !ok {
		return nil, fmt.Errorf(
			"%w:%s", ErrorStepTypeIsNotMatch, states.StateTypes.Suspend)
	}

	suspendState := &Suspend{
		ExecutionStep: baseStep,
		State:         state,
	}
	return suspendState, err
}

func (s *Suspend) Run(_ queue.InnerMessageBody) error {

	var err error
	var dbStep *po.Step
	var dbStepgroup *po.StepGroup

	dbStep = s.Data

	starttime := time.Now()

	tx, maker := s.ExecutionService.GetMetaDB().NewTxMaker(nil)
	defer maker.Close(&err)

	updatestate := po.Step{
		Status: string(StepStatus.Suspending),
	}

	err = tx.Where(po.Step{ID: dbStep.ID}).Updates(&updatestate).Error
	if err != nil {
		slog.Error("update step status error", "error", err)
		return err
	}
	tx.Commit()

	finishtime := time.Now()
	event := vo.ExecutionEvent{
		ExecutionID: dbStep.ExecutionID,
		StartTime:   starttime,
		FinishTime:  finishtime,
		StepID:      dbStep.ID,
		StepName:    dbStep.Name,
		Data:        EventContent_SuspendStepExecuted{},
	}
	s.ExecutionService.SendExecutionEvents(event)

	// message queue send create message
	message := NewStepMessage(dbStep.ExecutionID, MessageType.StepGroupSuspend, dbStepgroup.StepID, nil)
	err = s.ExecutionService.SendInnerMessage(message, nil)
	if err != nil {
		return err
	}

	return nil
}

// ProcessEventSuspendTimeout ProcessEventSuspendTimeout
func (s *Suspend) ProcessEventSuspendTimeout(message queue.InnerMessageBody) error {

	return nil
}
