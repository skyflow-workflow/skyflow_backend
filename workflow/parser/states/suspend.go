package states

// Suspend   suspend state
type Suspend struct {
	*BaseState
}

func (s *Suspend) GetBaseState() *BaseState {
	return s.BaseState
}
