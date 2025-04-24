package stepfunction

import (
	"context"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/decoder"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
)

// StepfuncionDecoder ...
type StepfuncionDecoder struct {
	*decoder.CommonDecoder
	config *decoder.ParserConfig
	quota  *decoder.Quota
}

// NewStepfuncionDecoder ...
func NewStepfuncionDecoder(config *decoder.ParserConfig, quota *decoder.Quota) *StepfuncionDecoder {
	return &StepfuncionDecoder{
		CommonDecoder: decoder.NewCommonDecoder(),
		config:        config,
		quota:         quota,
	}
}

// Decode ...
func (decoder *StepfuncionDecoder) Decode(definition string) (*states.StateMachine, error) {
	// Parse the state machine
	var err error
	datamap := make(map[string]any)
	err = decoder.JSONUnmarshall(definition, &datamap)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	// Check the type
	sm, err := decoder.DecodeStateMachine(ctx, datamap)
	if err != nil {
		return nil, err
	}
	return sm, nil
}
