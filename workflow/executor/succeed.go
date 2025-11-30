package executor

import (
	"fmt"
	"time"

	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backend/workflow/vo"
)

// Succeed executor for succeed state
type Succeed struct {
	*ExecutionStep
	State *states.FailState
}

// NewSucceedFromID NewSucceedFromID
func NewSucceedFromID(id int, executor *Executor) (*Succeed, error) {

	var err error

	var dbStep = &po.Step{}

	dbStep, err = executor.ExecutionService.QueryStepByID(id, []string{}, nil)
	if err != nil {
		return nil, err
	}

	state, err := NewSucceedFromData(dbStep, executor)

	return state, err
}

// NewSucceedFromData NewSucceedFromData
func NewSucceedFromData(dbStep *po.Step, executor *Executor) (*Succeed, error) {

	baseStep, err := NewExecutionStep(dbStep, executor)
	if err != nil {
		return nil, err
	}
	state, ok := baseStep.State.(*states.FailState)
	if !ok {
		return nil, fmt.Errorf(
			"%w:%s", ErrorStepTypeIsNotMatch, states.StateTypes.Succeed)
	}
	succeedState := &Succeed{
		ExecutionStep: baseStep,
		State:         state,
	}
	return succeedState, err
}

func (s *Succeed) Run(msgBody queue.InnerMessageBody) error {

	var err error
	var step_id = msgBody.StepID
	var execution_id = msgBody.ExecutionID

	starttime := time.Now()

	lock := s.Executor.ExecutionService.LockService.LockExecution(execution_id)
	defer lock.Unlock()

	lockTx := lock.GetTx()

	tx, maker := s.Executor.ExecutionService.GetMetaDB().NewTxMaker(lockTx)
	defer maker.Close(&err)

	updateStep := &po.Step{
		Status: string(StepStatus.Success),
	}
	updateExecution := &po.Execution{
		Status: string(ExecutionStatus.Success),
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

	event1 := vo.ExecutionEvent{
		ExecutionID: execution_id,
		StartTime:   starttime,
		FinishTime:  finishtime,
		StepID:      step_id,
		StepName:    s.Data.Name,
		Data: EventContent_SucceedStateExecuted{
			Output: s.Data.Output,
		},
	}
	event2 := vo.ExecutionEvent{
		ExecutionID: execution_id,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data: EventContent_ExecutionSucceeded{
			Output: s.Data.Output,
		},
	}
	s.Executor.ExecutionService.SendExecutionEvents(event1, event2)
	return nil
}
