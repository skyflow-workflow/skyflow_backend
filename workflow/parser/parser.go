// Description: This package contains the parser service and configuration for the workflow parser.
package parser

import (
	"encoding/json"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/decoder"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/stepfunction"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/quota"
)

// ParserConfig parser configuration
type Parser struct {
	Config decoder.ParserConfig
	Quota  quota.Quota
}

// ParserConfig parser configuration
func NewParser(config decoder.ParserConfig, quotaconfig quota.Quota) *Parser {
	return &Parser{
		Config: config,
		Quota:  quotaconfig,
	}
}

func ValdateStateMachine(definition string) error {
	// Validate the state machine
	return nil
}

func (parser *Parser) ParseStateMachine(definition string) (*states.StateMachine, error) {
	decoder := stepfunction.NewStepfuncionDecoder(&parser.Config, &parser.Quota)
	// Parse the state machine
	datamap := make(map[string]interface{})
	err := json.Unmarshal([]byte(definition), &datamap)
	if err != nil {
		return nil, err
	}
	sm, err := decoder.Decode(definition)
	if err != nil {
		return nil, err
	}
	return sm, nil
}
