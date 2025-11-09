// Description: This package contains the parser service and configuration for the workflow parser.
package parser

import (
	"context"
	"fmt"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/config"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/stepfunction"
)

// ParserConfig parser configuration
type Parser struct {
	Config              *config.Config
	stepfunctionDecoder *stepfunction.StepfuncionDecoder
}

// NewParser ParserConfig parser configuration
func NewParser(config *config.Config) *Parser {
	return &Parser{
		Config:              config,
		stepfunctionDecoder: stepfunction.NewStepfuncionDecoder(config),
	}
}

// ValdateStateMachine ...
func ValdateStateMachine(definition string) error {
	// Validate the state machine
	return nil
}

// ParseStateMachine ...
func (parser *Parser) ParseStateMachine(definition string) (*states.StateMachine, error) {
	sm, err := parser.stepfunctionDecoder.Decode(definition)
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
func (parser *Parser) ParseState(definition string) (states.State, error) {
	// Parse the step definition
	state, err := parser.stepfunctionDecoder.DecodeStateDefintion(context.Background(), definition)
	return state, err
}

func ParseStateMachine(definition string) (*states.StateMachine, error) {
	return StandardParser.ParseStateMachine(definition)
}

func ParseState(definition string) (states.State, error) {
	return StandardParser.ParseState(definition)
}

func GenerateActivityURI(namespace string, activityName string) string {
	return StandardParser.GenerateActivityURI(namespace, activityName)
}
func GenerateStateMachineURI(namespace string, stateMachineName string) string {
	return StandardParser.GenerateStateMachineURI(namespace, stateMachineName)
}
