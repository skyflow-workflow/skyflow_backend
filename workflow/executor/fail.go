package executor

import (
	"fmt"
	"time"

	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backend/workflow/vo"
)

// Fail executor for fail state
type Fail struct {
	*ExecutionStep
	State *states.FailState
}

// NewFailFromID NewFailFromID
func NewFailFromID(id int, executor *Executor) (*Fail, error) {

	var err error

	var dbStep = &po.Step{}

	dbStep, err = executor.ExecutionService.QueryStepByID(id, []string{}, nil)
	if err != nil {
		return nil, err
	}

	state, err := NewFailFromData(dbStep, executor)

	return state, err
}

// NewFailFromData NewFailFromData
func NewFailFromData(dbStep *po.Step, executor *Executor) (*Fail, error) {

	baseStep, err := NewExecutionStep(dbStep, executor)
	if err != nil {
		return nil, err
	}
	state, ok := baseStep.State.(*states.FailState)
	if !ok {
		return nil, fmt.Errorf(
			"%w:%s", ErrorStepTypeIsNotMatch, states.StateTypes.Fail)
	}
	failState := &Fail{
		ExecutionStep: baseStep,
		State:         state,
	}
	return failState, err
}

func (f *Fail) Run(msgBody queue.InnerMessageBody) error {

	var err error
	var step_id = msgBody.StepID
	var execution_id = msgBody.ExecutionID

	starttime := time.Now()

	lock := f.Executor.ExecutionService.LockService.LockExecution(execution_id)
	defer lock.Unlock()

	lockTx := lock.GetTx()

	tx, maker := f.Executor.ExecutionService.GetMetaDB().NewTxMaker(lockTx)
	defer maker.Close(&err)

	updateStep := &po.Step{
		Status: string(StepStatus.Failed),
	}
	updateExecution := &po.Execution{
		Status: string(ExecutionStatus.Failed),
	}
	err = tx.Where(po.Step{ID: step_id}).Updates(updateStep).Error
	if err != nil {
		return err
	}
	err = tx.Where(po.Execution{ID: execution_id}).Updates(updateExecution).Error
	if err != nil {
		return err
	}
	tx.Commit()

	lock.Unlock()

	finishtime := time.Now()

	// FailStateExecuted 事件， Step Event
	event1 := vo.ExecutionEvent{
		ExecutionID: execution_id,
		StartTime:   starttime,
		FinishTime:  finishtime,
		StepID:      step_id,
		StepName:    f.Data.Name,
		Data: EventContent_FailStateExecuted{
			Error: f.State.Error,
			Cause: f.State.Cause,
		},
	}
	// ExecutionFailed 事件， 在FailStateExecuted 事件之后发送，

	event2 := vo.ExecutionEvent{
		ExecutionID: execution_id,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data: EventContent_ExecutionFailed{
			EventType: msgBody.Type,
			Error:     f.State.Error,
			Cause:     f.State.Cause,
		},
	}
	f.Executor.ExecutionService.SendExecutionEvents(event1, event2)

	return nil
}
