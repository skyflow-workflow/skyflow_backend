package states

type StateMachineBody struct {
	States map[string]State
}

func (s *StateMachineBody) Validate() string {
	return "StateMachineBody"
}
