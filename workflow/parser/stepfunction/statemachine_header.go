package stepfunction

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/decoder"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
)

// DecodeStateMachineHeader ...
func (sfdecoder *StepfuncionDecoder) DecodeStateMachineHeader(ctx context.Context, data map[string]interface{}) (
	*states.StateMachineHeader, error) {
	// Parse the state machine
	var err error
	header := decoder.DefaultStateMachineHeader
	err = mapstructure.Decode(data, &header)
	if err != nil {
		return nil, decoder.MergeError(ctx, err)
	}
	err = header.Init()
	if err != nil {
		return nil, decoder.MergeError(ctx, err)
	}
	return &header, nil
}

// DecodeStateMachineHeaderDefintion ...
func (sfdecoder *StepfuncionDecoder) DecodeStateMachineHeaderDefintion(definition string) (*states.StateMachineHeader, error) {

	var err error
	datamap := make(map[string]interface{})
	err = sfdecoder.JSONUnmarshall(definition, &datamap)
	if err != nil {
		return nil, err
	}
	header, err := sfdecoder.DecodeStateMachineHeader(context.Background(), datamap)
	if err != nil {
		return nil, err
	}
	return header, nil
}
