// Decoder package
// This package defines the decoder interface for the statemachine.
// It is used to decode the statemachine definition into a StateMachine.

package decoder

import "github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"

// Decoder interface
type Decoder interface {
	Decode(definition string) (*states.StateMachine, error)
}
