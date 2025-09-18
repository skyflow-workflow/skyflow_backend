package states

type StateGroup struct {
	*BaseState
	*StateMachineBody
}

// DefaultState 默认初始化的 State 的value
var DefaultStateGroupState = BaseState{
	Name:            "",
	Type:            string(StateTypes.StateGroup),
	InputPath:       "$",
	OutputPath:      "$",
	Next:            "",
	ResultPath:      "",
	Parameters:      nil,
	MaxExecuteTimes: 1000,
	End:             true,
}

func (sg *StateGroup) Init() (err error) {
	err = sg.BaseState.Init()
	if err != nil {
		return
	}
	err = sg.StateMachineBody.Init()
	if err != nil {
		return
	}
	return nil
}

func (sg *StateGroup) GetBone() StateBone {
	basebone := sg.BaseState.GetBone()
	smbbone := sg.StateMachineBody.GetBone()
	basebone.StateMachineBone = &smbbone
	return basebone
}
