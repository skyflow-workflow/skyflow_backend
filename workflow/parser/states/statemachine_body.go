package states

import (
	"fmt"
)

// StateMachineBody ...
type StateMachineBody struct {
	StartAt string
	States  map[string]State
}

// AddState ...
func (s *StateMachineBody) AddState(state State) {
	s.States[state.GetName()] = state
}

func (s *StateMachineBody) SetStartAt(name string) {
	s.StartAt = name
}

func (s *StateMachineBody) GetBone() StateMachineBone {
	bone := StateMachineBone{
		StartAt: s.StartAt,
		States:  make(map[string]StateBone),
	}
	for name, state := range s.States {
		bone.States[name] = state.GetBone()
	}
	return bone
}

// Validate ...
func (s *StateMachineBody) Validate() error {

	// verify statemachine

	// verify StartAt
	if _, ok := s.States[s.StartAt]; !ok {
		return fmt.Errorf("field '%s' state '%s' not found in statemachine",
			StateMachineFieldNames.StartAt, s.StartAt)
	}
	// verify all nodes next in state valid
	for statename, statebone := range s.States {
		for _, next := range statebone.GetBone().Next {
			if _, ok := s.States[next]; !ok {
				return fmt.Errorf(
					"state '%s' Next '%s' not found in statemachine",
					statename, next,
				)
			}
		}
	}

	hasEnd := false
	for _, state := range s.States {
		if state.GetBone().End {
			hasEnd = true
			break
		}
	}
	if !hasEnd {
		return fmt.Errorf("statemachine need end state ")
	}
	return nil
}
