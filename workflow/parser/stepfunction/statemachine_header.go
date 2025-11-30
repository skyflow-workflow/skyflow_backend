package stepfunction

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/skyflow-workflow/skyflow_backend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/decoder"
	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"
)

// DecodeStateMachineHeader ...
func (sfdecoder *StepFunctionDecoder) DecodeStateMachineHeader(ctx context.Context, data map[string]interface{}) (
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
func (sfdecoder *StepFunctionDecoder) DecodeStateMachineHeaderDefinition(definition string) (*states.StateMachineHeader, error) {

	var err error
	datamap, err := toolkit.DecodeStringToMap(definition)
	if err != nil {
		return nil, err
	}
	header, err := sfdecoder.DecodeStateMachineHeader(context.Background(), datamap)
	if err != nil {
		return nil, err
	}
	return header, nil
}
