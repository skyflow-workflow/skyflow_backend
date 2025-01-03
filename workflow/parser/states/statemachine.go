package states

type StateMachine struct {
	*StateMachineHeader
	States map[string]State
}

func (s *StateMachine) Validate() string {
	return "StateMachine"
}

func ParserStateMachine(definition string) error {
	// Parse the state machine
	return nil
}
