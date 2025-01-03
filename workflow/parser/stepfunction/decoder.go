package stepfunction

import (
	"encoding/json"

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
		return nil, err
	}

	return nil, nil
}

func (decoder *StepfuncionDecoder) DecodeStateMachine(definition string) (*states.StateMachine, error) {
	// Parse the state machine
	datamap := make(map[string]interface{})
	err := json.Unmarshal([]byte(definition), &datamap)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (decoder *StepfuncionDecoder) DecodeState(definition string) (*states.State, error) {
	// Parse the state machine
	datamap := make(map[string]interface{})
	err := json.Unmarshal([]byte(definition), &datamap)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
