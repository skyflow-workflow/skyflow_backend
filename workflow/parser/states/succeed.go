package states

// SucceedState  失败节点
type SucceedState struct {
	*BaseState
}

// NewSucceedStateFromString NewSucceedStateFromString
func NewSucceedStateFromString(definition string) (state *SucceedState, err error) {

	data, err := StringToMap(definition)
	if err != nil {
		return nil, err
	}
	state, err = NewSucceedStateFromMap(data)
	return state, err
}

// NewSucceedStateFromMap  Create New Succeed State
func NewSucceedStateFromMap(data map[string]interface{}) (state *SucceedState, err error) {

	bs, err := NewBaseStateFromMap(data)
	if err != nil {
		return
	}
	bs.End = true
	state = &SucceedState{
		BaseState: bs,
	}
	return
}

// IsEnd 是否是终止节点
func (s *SucceedState) IsEnd() bool {
	return true
}

func (s *SucceedState) GetBaseState() *BaseState {
	return s.BaseState
}
