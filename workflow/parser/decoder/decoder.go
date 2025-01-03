package decoder

import "github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"

// Decoder interface
type Decoder interface {
	Decode(definition string) (*states.StateMachine, error)
}
