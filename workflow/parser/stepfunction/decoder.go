package stepfunction

import (
	"context"

	"github.com/skyflow-workflow/skyflow_backbend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/config"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/decoder"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
)

// StepFunctionDecoder decodes AWS Step Functions JSON definitions
type StepFunctionDecoder struct {
	*decoder.CommonDecoder
	config *config.Config
}

// NewStepFunctionDecoder creates a new StepFunctionDecoder instance
func NewStepFunctionDecoder(config *config.Config) *StepFunctionDecoder {
	return &StepFunctionDecoder{
		CommonDecoder: decoder.NewCommonDecoder(),
		config:        config,
	}
}

// Decode decodes a state machine definition from JSON string
func (decoder *StepFunctionDecoder) Decode(definition string) (*states.StateMachine, error) {
	// Parse the state machine
	var err error
	datamap, err := toolkit.DecodeStringToMap(definition)
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
