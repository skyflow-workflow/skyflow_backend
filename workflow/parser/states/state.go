// Package states implements the State behavior for the workflow
// states are the basic step in the workflow, they can be combined to form a complex workflow
package states

import (
	"fmt"
	"time"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

// State ...
type State interface {
	Validate() error
	GetName() string
	SetName(name string)
	GetType() string
	GetBone() StateBone
	IsEnd() bool
	GetNext() []string
	GetBaseState() *BaseState
	GetDefinition() (map[string]any, error)
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

	typeKey, ok := data[StateFieldNames.Type]
	if !ok {
		err = fmt.Errorf(" state lack of field [ Type ] ")
		return state, err
	}
	typeKeyStr, ok := typeKey.(string)
	if !ok {
		err = fmt.Errorf(" field [ Type ] field is not string  ")
		return
	}
	switch StateType(typeKeyStr) {
	case StateTypes.Task:
		state, err = NewTaskStateFromMap(data)
	case StateTypes.Choice:
		state, err = NewChoiceStateFromMap(data)
	case StateTypes.Wait:
		state, err = NewWaitStateFromMap(data)
	case StateTypes.Pass:
		state, err = NewPassStateFromMap(data)
	case StateTypes.Fail:
		state, err = NewFailStateFromMap(data)
	case StateTypes.Succeed:
		state, err = NewSucceedStateFromMap(data)
	case StateTypes.Parallel:
		state, err = NewParallelStateFromMap(data, depth)
	// case StateTypes.Map:
	// 	state, err = NewMapStateFromMap(data, depth)

	default:
		err = fmt.Errorf("%w : %s", vo.ErrorUnsupportedStateType, typeKeyStr)
		return
	}

	if err != nil {
		return nil, err
	}
	return

}
