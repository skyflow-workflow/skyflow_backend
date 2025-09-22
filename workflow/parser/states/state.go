// Package states implements the State behavior for the workflow
// states are the basic step in the workflow, they can be combined to form a complex workflow
package states

import (
	"fmt"
	"time"
)

// State ...
type State interface {
	Init() error
	Validate() error
	GetName() string
	SetName(name string)
	GetType() string
	GetBone() StateBone
	IsEnd() bool
	GetNext() []string
	GetBaseState() *BaseState
}

// NextState next state info
type NextState struct {
	Name       string        // Next State Name
	Output     any           // State Output as the input of the next state
	Delay      time.Duration // Delay delay duration from current state to next state
	Retry      bool          // Whether to trigger retry
	RetryIndex int           // Which Indexed Retry strategy hit
}

// NewStateFromMap  NewStateFromMap
func NewStateFromMap(data map[string]interface{}, depth int) (state State, err error) {

	typeKey, ok := data[Fields.Type]
	if !ok {
		err = fmt.Errorf(" state lack of field [ Type ] ")
		return state, err
	}
	typeKeyStr, ok := typeKey.(string)
	if !ok {
		err = fmt.Errorf(" field [ Type ] field is not string  ")
		return
	}
	switch typeKeyStr {
	case StateType.Task:
		state, err = NewTaskStateFromMap(data)
	case StateType.Choice:
		state, err = NewChoiceStateFromMap(data)
	case StateType.Wait:
		state, err = NewWaitStateFromMap(data)
	case StateType.Pass:
		state, err = NewPassStateFromMap(data)
	case StateType.Fail:
		state, err = NewFailStateFromMap(data)
	case StateType.Succeed:
		state, err = NewSucceedStateFromMap(data)
	case StateType.Parallel:
		state, err = NewParallelStateFromMap(data, depth)
	case StateType.Map:
		state, err = NewMapStateFromMap(data, depth)
	case StateType.StateGroup:
		state, err = NewStateGroupFromMap(data, depth)
	default:
		err = fmt.Errorf("%w : %s", vo.ErrorUnrecognizeStatemachineType, typeKeyStr)
		return
	}
	if err != nil {
		return
	}
	err = state.Init()
	if err != nil {
		err = fmt.Errorf("step '%s'  error : %w", state.GetName(), err)
		return
	}
	return

}
