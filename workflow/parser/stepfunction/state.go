package stepfunction

import (
	"context"
	"fmt"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/decoder"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
)

// DecodeBaseState ...
func (sfDecoder *StepfuncionDecoder) DecodeBaseState(ctx context.Context, data map[string]any) (states.State, error) {

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

// DecodeStateDefintion ...
func (sfDecoder *StepfuncionDecoder) DecodeStateDefintion(ctx context.Context, definition string) (states.State, error) {
	var err error
	datamap := make(map[string]any)
	err = sfDecoder.JSONUnmarshal(definition, &datamap)
	if err != nil {
		return nil, err
	}
	state, err := sfDecoder.DecodeState(ctx, datamap)
	if err != nil {
		return nil, err
	}
	return state, nil

}

// DecodeState ...
func (decoder *StepfuncionDecoder) DecodeState(ctx context.Context, data map[string]any) (states.State, error) {

	state, err := states.NewStateFromMap(data, states.StartDepth)
	if err != nil {
		return nil, err
	}
	return state, nil
}
