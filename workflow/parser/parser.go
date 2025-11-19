// Description: This package contains the parser service and configuration for the workflow parser.
package parser

import (
	"context"
	"fmt"

	"github.com/skyflow-workflow/skyflow_backend/workflow/config"
	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/stepfunction"
)

// Parser parser configuration for workflow definitions
type Parser struct {
	Config              *config.Config
	stepfunctionDecoder *stepfunction.StepFunctionDecoder
}

// NewParser creates a new Parser instance with the given configuration
func NewParser(config *config.Config) *Parser {
	return &Parser{
		Config:              config,
		stepfunctionDecoder: stepfunction.NewStepFunctionDecoder(config),
	}
}

// ValidateStateMachine validates a state machine definition
func ValidateStateMachine(definition string) error {
	if definition == "" {
		return fmt.Errorf("definition cannot be empty")
	}

	if StandardParser == nil {
		return fmt.Errorf("StandardParser is not initialized")
	}

	// Validate the state machine
	_, err := StandardParser.ParseStateMachine(definition)
	if err != nil {
		return err
	}
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
	if definition == "" {
		return nil, fmt.Errorf("definition cannot be empty")
	}

	// Parse the step definition
	state, err := parser.stepfunctionDecoder.DecodeStateDefinition(context.Background(), definition)
	return state, err
}

// ParseStateMachine parses a state machine definition using the standard parser
func ParseStateMachine(definition string) (*states.StateMachine, error) {
	if definition == "" {
		return nil, fmt.Errorf("definition cannot be empty")
	}

	if StandardParser == nil {
		return nil, fmt.Errorf("StandardParser is not initialized")
	}

	return StandardParser.ParseStateMachine(definition)
}

// ParseState parses a state definition using the standard parser
func ParseState(definition string) (states.State, error) {
	if definition == "" {
		return nil, fmt.Errorf("definition cannot be empty")
	}

	if StandardParser == nil {
		return nil, fmt.Errorf("StandardParser is not initialized")
	}

	return StandardParser.ParseState(definition)
}

// GenerateActivityURI generates an activity URI using the standard parser
func GenerateActivityURI(namespace string, activityName string) string {
	if StandardParser == nil {
		return ""
	}

	return StandardParser.GenerateActivityURI(namespace, activityName)
}

// GenerateStateMachineURI generates a state machine URI using the standard parser
func GenerateStateMachineURI(namespace string, stateMachineName string) string {
	if StandardParser == nil {
		return ""
	}

	return StandardParser.GenerateStateMachineURI(namespace, stateMachineName)
}
