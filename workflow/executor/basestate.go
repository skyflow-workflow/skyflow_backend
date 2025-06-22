package executor

import (
	"fmt"

	"github.com/skyflow-workflow/skyflow_backbend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
)

// ExecutionStep execution step
type ExecutionStep struct {
	Data      *po.Step
	Executor  *Executor
	State     states.State
	BaseState *states.BaseState
}

// NewBaseExecutionStep New一个可执行状态实例
func NewExecutionStep(dbStep *po.Step, executor *Executor) (*ExecutionStep, error) {
	var err error

	state, err := executor.parser.ParseState(dbStep.Definition)
	if err != nil {
		return nil, err
	}
	exeStep := &ExecutionStep{
		State:     state,
		Executor:  executor,
		Data:      dbStep,
		BaseState: state.GetBaseState(),
	}
	return exeStep, nil

}

// GetBone  Get Execution Step Status and return StepBone  struct
func (s *ExecutionStep) GetBone() StepBone {

	bone := StepBone{
		BaseBone: BaseBone{
			BaseBone: s.State.GetBone().BaseBone,
			Status:   s.Data.Status,
			StepID:   s.Data.ID,
		},
	}
	if s.Data.ExecuteIndex > 0 {
		bone.Index = s.Data.ExecuteIndex
	}
	return bone
}

// ProcessEvent process step event
// ProcessEvent is a placeholder method for processing events.
// It currently returns an error indicating that the method is not implemented.
func (s *ExecutionStep) ProcessEvent(message queue.InnerMessage) error {
	return fmt.Errorf("method not implement")
}

// DecodeInput decode intput to json
func (s *ExecutionStep) DecodeInput() (any, error) {
	i, err := toolkit.DecodeStringToMap(s.Data.Input)
	return i, err
}

// DecodeOutput decode output to json
func (s *ExecutionStep) DecodeOutput() (any, error) {
	i, err := toolkit.DecodeStringToMap(s.Data.Output)
	return i, err
}

// GetInput 获得经过计算的input
func (s *ExecutionStep) GetInput() (any, error) {
	i, err := s.DecodeInput()
	if err != nil {
		return nil, err
	}
	ii, err := s.State.GetBaseState().GetParametersInput(i)
	return ii, err
}

// GetNextStep  return (NextState string , output any, err error )
func (s *ExecutionStep) GetNextStep(output any) (NextStep, error) {

	var err error
	var sns NextStep
	var input any
	input, err = s.DecodeInput()
	if err != nil {
		return sns, nil
	}
	ns, err := s.State.GetBaseState().GetNextState(input, output)
	if err != nil {
		return sns, err
	}
	nextstep := NextStep{
		NextState: ns,
		GroupID:   s.Data.GroupID,
	}
	return nextstep, nil

}
