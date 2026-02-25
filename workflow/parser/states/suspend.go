package states

// SuspendState   suspend state
type SuspendState struct {
	*BaseState `mapstructure:"BaseState"`
}

func (s *SuspendState) GetBaseState() *BaseState {
	return s.BaseState
}

// NewSuspendStateFromString NewSuspendStateFromString
func NewSuspendStateFromString(definition string) (state *SuspendState, err error) {

	data, err := StringToMap(definition)
	if err != nil {
		return nil, err
	}
	state, err = NewSuspendStateFromMap(data)
	return state, err
}

// NewSucceedStateFromMap  Create New Succeed State
func NewSuspendStateFromMap(data map[string]interface{}) (state *SuspendState, err error) {

	bs, err := NewBaseStateFromMap(data)
	if err != nil {
		return
	}

	state = &SuspendState{
		BaseState: bs,
	}
	return
}

func (s *SuspendState) GetDefinition() (string, error) {
	baseData, err := s.BaseState.GetDefinitionMap()
	if err != nil {
		return "", err
	}
	dataStr, err := ToString(baseData)
	return dataStr, err
}
