package states

// Succeed  失败节点
type Succeed struct {
	*BaseState
}

// IsEnd 是否是终止节点
func (s *Succeed) IsEnd() bool {
	return true
}

func (s *Succeed) GetBaseState() *BaseState {
	return s.BaseState
}
