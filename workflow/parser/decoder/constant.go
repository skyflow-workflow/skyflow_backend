package decoder

import "github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"

// DefaultStateMachineHeader ...
var DefaultStateMachineHeader = states.StateMachineHeader{
	Version: "1.0",
	Type:    WorkflowType.Statemachine,
	// QueryLanguage: string(states.QueryLanguages.JSONPath),
}

// WorkflowType ...
var WorkflowType = struct {
	Pipeline     string
	Statemachine string
	Stepfunction string
}{
	Pipeline:     "pipeline",
	Statemachine: "statemachine",
	Stepfunction: "stepfunction",
}
