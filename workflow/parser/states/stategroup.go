package states

import (
	"encoding/json"
	"fmt"
)

type StateGroup struct {
	*BaseState
	StateMachineBody
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

// NewStateGroupFromString create statemachine from string
func NewStateGroupFromString(definition string, depth int) (sm *StateGroup, err error) {
	var jsonMap map[string]interface{}

	err = json.Unmarshal([]byte(definition), &jsonMap)
	if err != nil {
		return nil, err
	}
	sm, err = NewStateGroupFromMap(jsonMap, depth)
	return
}

// NewStateGroupFromMap Parse data to StateMachine
func NewStateGroupFromMap(data map[string]interface{}, depth int) (sg *StateGroup, err error) {

	// 状态机深度不能超过最大深度
	if depth > MaxDepth {
		err = fmt.Errorf("statemachine depth  reach MaxDepth [ %d ] ", MaxDepth)
		return nil, err
	}
	sg = NewStateGroup(depth)

	err = sg.InitByMap(data)
	if err != nil {
		return
	}
	err = sg.SetTypeDefinition(sg)
	if err != nil {
		return
	}
	return
}

func NewStateGroup(depth int) *StateGroup {

	bs := DefaultStateGroupState
	sg := &StateGroup{
		BaseState:        &bs,
		StateMachineBody: NewStateMachineBody(depth),
	}
	return sg
}

func (sg *StateGroup) InitByMap(data map[string]interface{}) (err error) {
	err = sg.BaseState.InitByMap(data)
	if err != nil {
		return
	}

	err = sg.StateMachineBody.InitByMap(data)
	if err != nil {
		return
	}
	return nil
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
