package stepfunction

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/decoder"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/quota"
)

type StepfuncionDecoder struct {
	config *decoder.ParserConfig
	quota  *quota.Quota
}

func NewStepfuncionDecoder(config *decoder.ParserConfig, quota *quota.Quota) *StepfuncionDecoder {
	return &StepfuncionDecoder{
		config: config,
		quota:  quota,
	}
}

func (decoder *StepfuncionDecoder) Decode(definition string) (*states.StateMachine, error) {
	// Parse the state machine
	datamap := make(map[string]interface{})
	err := json.Unmarshal([]byte(definition), &datamap)
	if err != nil {
		if jsonerr, ok := err.(*json.SyntaxError); ok {
			newerr := states.FieldError{
				RawError: err,
				Offset:   jsonerr.Offset,
			}
			return nil, &newerr
		}
		newerr := states.NewFieldError(err)
		return nil, newerr
	}

	// Check the type
	sm, err := decoder.DecodeStateMachine(datamap)
	if err != nil {
		return nil, err
	}
	return sm, nil
}

func (decoder *StepfuncionDecoder) DecodeStateMachine(data map[string]interface{}) (*states.StateMachine, error) {
	// Parse the state machine
	var err error
	header := states.DefaultStateMachineHeader
	err = mapstructure.Decode(data, &header)
	if err != nil {
		return nil, err
	}
	err = header.Init()
	if err != nil {
		return nil, err
	}

	body := StateMachineBody{}
	err = mapstructure.Decode(data, &body)
	if err != nil {
		return nil, err
	}
	smbody := states.StateMachineBody{
		States: make(map[string]states.State, len(body.States)),
	}

	for name, state := range body.States {
		newstate, err := decoder.DecodeState(state)
		if err != nil {
			return nil, err
		}
		newstate.SetName(name)
		smbody.States[name] = newstate
	}
	sm := &states.StateMachine{
		StateMachineHeader: &header,
		StateMachineBody:   &smbody,
	}
	return sm, nil
}

func (decoder *StepfuncionDecoder) DecodeState(data map[string]interface{}) (states.State, error) {

	var err error
	basestate := states.BaseState{}
	err = mapstructure.Decode(data, &basestate)
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
		state, err = decoder.DecodeTaskState(&basestate, data)
	default:
		err = states.NewFieldError(states.ErrInvalidStateType)
	}
	if err != nil {
		return nil, err
	}
	return state, err
}

func (decoder *StepfuncionDecoder) DecodeTaskState(basestate *states.BaseState, data map[string]interface{}) (
	states.State, error) {
	var err error
	taskbody := states.TaskBody{}
	err = mapstructure.Decode(data, &taskbody)
	if err != nil {
		return nil, err
	}
	taskstate := &states.Task{
		BaseState: basestate,
		TaskBody:  taskbody,
	}
	err = taskstate.Init()
	if err != nil {
		return nil, err
	}
	return taskstate, nil
}

func (decoder *StepfuncionDecoder) DecodeChoiceState(basestate *states.BaseState, data map[string]interface{}) (
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

func (decoder *StepfuncionDecoder) DecodeWaitState(basestate *states.BaseState, data map[string]interface{}) (
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
