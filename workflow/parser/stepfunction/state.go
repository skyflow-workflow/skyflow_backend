package stepfunction

import (
	"context"
	"fmt"

	"github.com/skyflow-workflow/skyflow_backend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/decoder"
	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"
)

// DecodeBaseState decodes a base state from the given data map
func (sfDecoder *StepFunctionDecoder) DecodeBaseState(ctx context.Context, data map[string]any) (states.State, error) {

	var err error
	basestate := states.BaseState{}
	err = sfDecoder.MapDecode(data, &basestate)
	if err != nil {
		return nil, err
	}
	err = basestate.Init()
	if err != nil {
		return nil, err
	}
	var state states.State
	switch basestate.Type {
	case string(states.StateTypes.Task):
		state, err = sfDecoder.DecodeTaskState(ctx, &basestate, data)
	default:
		curpath := append(decoder.GetPath(ctx), states.StateFieldNames.Type)
		err = states.NewFieldPathError(
			fmt.Errorf("%w: %s", states.ErrorInvalidStateType, basestate.Type),
			curpath...)
	}
	if err != nil {
		return nil, err
	}
	return state, err
}

// DecodeStateDefinition decodes a state definition from JSON string
func (sfDecoder *StepFunctionDecoder) DecodeStateDefinition(ctx context.Context, definition string) (states.State, error) {
	var err error
	datamap, err := toolkit.DecodeStringToMap(definition)
	if err != nil {
		return nil, err
	}
	state, err := sfDecoder.DecodeState(ctx, datamap)
	if err != nil {
		return nil, err
	}
	return state, nil

}

// DecodeState decodes a state from the given data map
func (decoder *StepFunctionDecoder) DecodeState(ctx context.Context, data map[string]any) (states.State, error) {

	state, err := states.NewStateFromMap(data, states.StartDepth)
	if err != nil {
		return nil, err
	}
	return state, nil
}
