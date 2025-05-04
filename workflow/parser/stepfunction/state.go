package stepfunction

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
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
func (sfDecoder *StepfuncionDecoder) DecodeStateDefintion(definition string) (states.State, error) {
	var err error
	datamap := make(map[string]any)
	err = sfDecoder.JSONUnmarshall(definition, &datamap)
	if err != nil {
		return nil, err
	}
	state, err := sfDecoder.DecodeState(context.Background(), datamap)
	if err != nil {
		return nil, err
	}
	return state, nil

}

// DecodeState ...
func (decoder *StepfuncionDecoder) DecodeState(ctx context.Context, data map[string]any) (states.State, error) {

	var err error

	err = states.ValidateStateFieldOptional(data)
	if err != nil {
		return nil, err
	}

	basestate := DefaultBaseState
	err = decoder.MapDecode(data, &basestate)
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
		state, err = decoder.DecodeTaskState(ctx, &basestate, data)
	case string(states.StateTypes.Choice):
		state, err = decoder.DecodeChoiceState(ctx, &basestate, data)
	case string(states.StateTypes.Pass):
		state, err = decoder.DecodePassState(ctx, &basestate, data)
	case string(states.StateTypes.Wait):
		state, err = decoder.DecodeWaitState(ctx, &basestate, data)
	default:
		rawerr := fmt.Errorf("%w: %s", states.ErrorInvalidStateType, basestate.Type)
		ctx = decoder.AddCtxDecodePath(ctx, states.StateFieldNames.Type)
		err = decoder.NewFieldPathError(ctx, rawerr)
	}
	if err != nil {
		return nil, err
	}
	return state, err
}

// DecodeWaitState ...
func (sfDecoder *StepfuncionDecoder) DecodeWaitState(ctx context.Context,
	basestate *states.BaseState, data map[string]any) (
	states.State, error) {
	var err error
	waitbody := &states.WaitBody{}
	err = mapstructure.Decode(data, waitbody)
	if err != nil {
		return nil, err
	}
	waitstate := &states.Wait{
		BaseState: basestate,
		WaitBody:  waitbody,
	}
	err = waitstate.Init()
	if err != nil {
		return nil, err
	}
	return waitstate, nil
}

// DecodePassState ...
func (sfDecoder *StepfuncionDecoder) DecodePassState(ctx context.Context,
	basestate *states.BaseState, data map[string]any) (
	states.State, error) {
	var err error
	passbody := &states.PassBody{}
	err = mapstructure.Decode(data, &passbody)
	if err != nil {
		return nil, err
	}
	passstate := &states.Pass{
		BaseState: basestate,
		PassBody:  passbody,
	}
	err = passstate.Validate()
	if err != nil {
		return nil, err
	}
	return passstate, nil
}
