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

func NewStateFromMap(data map[string]interface{}, depth int) (State, error) {
	var err error
	var state State
	switch data["Type"] {
	case string(StateTypes.Task):
		state = &Task{}
	case string(StateTypes.Choice):
		state = &Choice{}
	case string(StateTypes.Pass):
		state = &Pass{}
	case string(StateTypes.Suspend):
		state = &Suspend{}
	case string(StateTypes.Wait):
		state = &Wait{}
	case string(StateTypes.Fail):
		state = &Fail{}
	case string(StateTypes.Succeed):
		state = &Succeed{}
	default:
		err = fmt.Errorf("state [ %s ] type : [ %s ] not supported", data["Name"], data["Type"])
	}
	if err != nil {
		return nil, err
	}
	return state, nil
}
