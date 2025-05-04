// Description: This package contains the parser service and configuration for the workflow parser.
package parser

import (
	"fmt"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/decoder"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/stepfunction"
)

// ParserConfig parser configuration
type Parser struct {
	Config decoder.ParserConfig
	Quota  decoder.Quota
}

// NewParser ParserConfig parser configuration
func NewParser(config decoder.ParserConfig, quotaconfig decoder.Quota) *Parser {
	return &Parser{
		Config: config,
		Quota:  quotaconfig,
	}
}

// ValdateStateMachine ...
func ValdateStateMachine(definition string) error {
	// Validate the state machine
	return nil
}

// ParseStateMachine ...
func (parser *Parser) ParseStateMachine(definition string) (*states.StateMachine, error) {
	decoder := stepfunction.NewStepfuncionDecoder(&parser.Config, &parser.Quota)
	sm, err := decoder.Decode(definition)
	if err != nil {
		return nil, err
	}
	return sm, nil
}
func (parser *Parser) GenerateActivityURI(namespace string, activityName string) string {
	// Generate the activity URI
	activity_uri := fmt.Sprintf("%s:%s/%s", "activity", namespace, activityName)
	return activity_uri
}
func (parser *Parser) GenerateStateMachineURI(namespace string, stateMachineName string) string {
	// Generate the workflow URI
	workflow_uri := fmt.Sprintf("%s:%s/%s", "statemachine", namespace, stateMachineName)
	return workflow_uri
}

func ParseStateMachine(definition string) (*states.StateMachine, error) {
	return StandardParser.ParseStateMachine(definition)
}

func GenerateActivityURI(namespace string, activityName string) string {
	return StandardParser.GenerateActivityURI(namespace, activityName)
}
func GenerateStateMachineURI(namespace string, stateMachineName string) string {
	return StandardParser.GenerateStateMachineURI(namespace, stateMachineName)
}
