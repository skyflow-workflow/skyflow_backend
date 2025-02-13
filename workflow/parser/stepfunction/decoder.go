package stepfunction

import (
	"context"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/decoder"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/quota"
)

// StepfuncionDecoder ...
type StepfuncionDecoder struct {
	*decoder.CommonDecoder
	config *decoder.ParserConfig
	quota  *quota.Quota
}

// NewStepfuncionDecoder ...
func NewStepfuncionDecoder(config *decoder.ParserConfig, quota *quota.Quota) *StepfuncionDecoder {
	return &StepfuncionDecoder{
		CommonDecoder: decoder.NewCommonDecoder(),
		config:        config,
		quota:         quota,
	}
}

// Decode ...
func (sfdecoder *StepfuncionDecoder) Decode(definition string) (*states.StateMachine, error) {
	// Parse the state machine
	var err error
	datamap := make(map[string]interface{})
	err = sfdecoder.JSONUnmashall(definition, &datamap)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	// Check the type
	sm, err := sfdecoder.DecodeStateMachine(ctx, datamap)
	if err != nil {
		return nil, err
	}
	return sm, nil
}

// DecodeStateMachine ...
func (sfdecoder *StepfuncionDecoder) DecodeStateMachine(ctx context.Context, data map[string]interface{}) (*states.StateMachine, error) {
	// Parse the state machine
	var err error
	var header *states.StateMachineHeader
	var body *states.StateMachineBody
	header, err = sfdecoder.DecodeStateMachineHeader(ctx, data)
	if err != nil {
		return nil, err
	}

	body, err = sfdecoder.DecodeStateMachineBody(ctx, data)
	if err != nil {
		return nil, err
	}
	sm := &states.StateMachine{
		StateMachineHeader: header,
		StateMachineBody:   body,
	}
	return sm, nil
}

// DecodeStateMachineBody ...
func (sfdecoder *StepfuncionDecoder) DecodeStateMachineBody(ctx context.Context, data map[string]interface{}) (
	*states.StateMachineBody, error) {

	var err error
	bodydata := StateMachineBody{}
	err = sfdecoder.MapDecode(data, &bodydata)
	if err != nil {
		return nil, err
	}
	smbody := &states.StateMachineBody{
		States: make(map[string]states.State, len(bodydata.States)),
	}

	ctx = decoder.AddPath(ctx, "States")
	for name, state := range bodydata.States {
		statectx := decoder.AddPath(ctx, name)
		newstate, err := sfdecoder.DecodeState(statectx, state)
		if err != nil {
			return nil, err
		}
		newstate.SetName(name)
		smbody.States[name] = newstate
	}

	return smbody, nil
}
