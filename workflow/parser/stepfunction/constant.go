package stepfunction

import "github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"

// StateMachineFields ...
var StateMachineFields = struct {
	Comment string
	Version string
	Type    string
	States  string
}{
	Comment: "Comment",
	Version: "Version",
	Type:    "Type",
	States:  "States",
}

// DefaultBaseState default basestate
var DefaultBaseState = states.BaseState{
	MaxExecuteTimes: 1000,
	InputPath:       "$",
	OutputPath:      "$",
	ResultPath:      "",
}
