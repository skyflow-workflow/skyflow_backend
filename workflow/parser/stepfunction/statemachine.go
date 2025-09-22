package stepfunction

import (
	"context"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
)

// StateMachineBody ...
type StateMachineBody struct {
	StartAt string
	States  map[string]map[string]any
}

// DecodeStateMachine ...
func (decoder *StepfuncionDecoder) DecodeStateMachine(ctx context.Context, data map[string]any) (*states.StateMachine, error) {

	sm, err := states.NewStateMachineFromMap(data)
	if err != nil {
		return nil, err
	}
	return sm, nil
}
