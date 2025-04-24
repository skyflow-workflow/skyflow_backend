// Package states implements the State behavior for the workflow
// states are the basic step in the workflow, they can be combined to form a complex workflow
package states

import "time"

// State ...
type State interface {
	Init() error
	Validate() error
	GetName() string
	SetName(name string)
	GetType() string
	GetBone() StateBone
}

// NextState next state info
type NextState struct {
	Name       string        // Next State Name
	Output     interface{}   // State Output as the input of the next state
	Delay      time.Duration // Delay delay duration from current state to next state
	Retry      bool          // Whether to trigger retry
	RetryIndex int           // Which Indexed Retry strategy hit
}
