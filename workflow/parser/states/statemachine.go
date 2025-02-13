package states

// StateMachine ...
type StateMachine struct {
	*StateMachineHeader
	*StateMachineBody
}

// Validate ...
func (s *StateMachine) Validate() string {
	return "StateMachine"
}

// ParserStateMachine ...
func ParserStateMachine(definition string) error {
	// Parse the state machine
	return nil
}
