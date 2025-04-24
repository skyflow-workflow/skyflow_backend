package stepfunction

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/decoder"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
)

// DecodeBaseState ...
func (sfdecoder *StepfuncionDecoder) DecodeBaseState(ctx context.Context, data map[string]interface{}) (states.State, error) {

	var err error
	basestate := states.BaseState{}
	err = sfdecoder.MapDecode(data, &basestate)
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
		state, err = sfdecoder.DecodeTaskState(&basestate, data)
	default:
		curpath := append(decoder.GetPath(ctx), states.StateFields.Type)
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
func (sfdecoder *StepfuncionDecoder) DecodeStateDefintion(definition string) (states.State, error) {
	var err error
	datamap := make(map[string]interface{})
	err = sfdecoder.JSONUnmashall(definition, &datamap)
	if err != nil {
		return nil, err
	}
	state, err := sfdecoder.DecodeState(context.Background(), datamap)
	if err != nil {
		return nil, err
	}
	return state, nil

}

// DecodeState ...
func (sfdecoder *StepfuncionDecoder) DecodeState(ctx context.Context, data map[string]interface{}) (states.State, error) {

	var err error
	// validate field requirment
	fieldrequired := states.BaseStateFieldRequired{}

	err = sfdecoder.MapDecode(data, &fieldrequired)
	if err != nil {
		return nil, err
	}
	err = fieldrequired.Validate()
	if err != nil {
		return nil, err
	}

	basestate := DefaultBaseState
	err = sfdecoder.MapDecode(data, &basestate)
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
		state, err = sfdecoder.DecodeTaskState(&basestate, data)
	case string(states.StateTypes.Choice):
		state, err = sfdecoder.DecodeChoiceState(&basestate, data)
	case string(states.StateTypes.Pass):
		state, err = sfdecoder.DecodePassState(&basestate, data)
	default:
		err = states.NewFieldPathError(fmt.Errorf("%w: %s", states.ErrorInvalidStateType, basestate.Type), states.StateFields.Type)
	}
	if err != nil {
		return nil, err
	}
	return state, err
}

// DecodeChoiceState ...
func (sfdecoder *StepfuncionDecoder) DecodeChoiceState(basestate *states.BaseState, data map[string]interface{}) (
	states.State, error) {
	var err error
	choicelbody := states.ChoiceBody{}
	err = mapstructure.Decode(data, &choicelbody)
	if err != nil {
		return nil, err
	}
	choicestate := &states.Choice{
		BaseState:  basestate,
		ChoiceBody: choicelbody,
	}
	err = choicestate.Init()
	if err != nil {
		return nil, err
	}
	return choicestate, nil
}

// DecodeWaitState ...
func (sfdecoder *StepfuncionDecoder) DecodeWaitState(basestate *states.BaseState, data map[string]interface{}) (
	states.State, error) {
	var err error
	waitbody := states.WaitBody{}
	err = mapstructure.Decode(data, &waitbody)
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
func (sfdecoder *StepfuncionDecoder) DecodePassState(basestate *states.BaseState, data map[string]interface{}) (
	states.State, error) {
	var err error
	waitbody := states.WaitBody{}
	err = mapstructure.Decode(data, &waitbody)
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
