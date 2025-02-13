package states

// StateMachineBody ...
type StateMachineBody struct {
	States map[string]State
}

// Validate ...
func (s *StateMachineBody) Validate() string {
	return "StateMachineBody"
}
