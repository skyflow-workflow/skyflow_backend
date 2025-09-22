package states

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// StateMachineGroupState statemachine
type StateMachineGroupState struct {
	GroupID    int            //自身的所处的group
	GroupIndex int            // 自身所在的GroupIndex
	Depth      int            //节点所在的 深度
	Name       string         // StateName
	State      State          //State
	SubGroup   *SubGroupState //子group列表
}

// SubGroupState statemachine中的子组
type SubGroupState struct {
	MasterGroupID int
	MasterName    string
	SubGroupID    int    `json:"group_id"`
	SubStartAt    string `json:"start_at"`
}

// StateMachineBody statemachine 的body 定义
type StateMachineBody struct {
	StartAt       string         `mapstructure:"StartAt" validate:"required,gt=0"`
	States        map[string]any `mapstructure:"States" validate:"required,gt=0"`
	_States       map[string]State
	_Depth        int
	_NewStateFunc func(data map[string]any, depth int) (State, error)
}

func (smb *StateMachineBody) SetNewStateFunc(f func(data map[string]any, depth int) (State, error)) {
	smb._NewStateFunc = f
}

func NewStateMachineBodyFromString(definition string, depth int) (smb *StateMachineBody, err error) {

	datamap, err := ToMap(definition)
	if err != nil {
		return
	}
	smb, err = NewStateMachineBodyFromMap(datamap, depth)
	if err != nil {
		return
	}
	return
}

// NewStateMachineBodyFromMap Parse data to StateMachine
func NewStateMachineBodyFromMap(data map[string]any, depth int) (smb *StateMachineBody, err error) {

	// 状态机深度不能超过最大深度
	if depth > MaxDepth {
		err = fmt.Errorf("statemachine depth  reach MaxDepth [ %d ] ", MaxDepth)
		return
	}
	newsmb := NewDefaultStateMachineBody(depth)
	smb = &newsmb

	err = InitStateMachineBodyByMap(smb, data)
	if err != nil {
		return
	}
	return
}

// InitByMap Use map to init StateMachineBody content
func InitStateMachineBodyByMap(smb *StateMachineBody, data map[string]any) (err error) {

	// 解析验证 剩余字段states 字段
	err = mapstructure.Decode(data, smb)
	if err != nil {
		return
	}

	err = myvalidate.Struct(smb)
	if err != nil {
		return
	}

	for name, stateObj := range smb.States {
		if name == "" {
			err = fmt.Errorf(" state name should not be '' ")
			return
		}
		content, ok := stateObj.(map[string]any)
		if !ok {
			err = fmt.Errorf(" state [ %s ] content should be map  ", name)
			return
		}
		var newstate State
		newstate, err = smb._NewStateFunc(content, smb._Depth)
		if err != nil {
			err = fmt.Errorf("state [ %s ] failed: %w", name, err)
			return
		}
		newstate.SetName(name)
		smb.AddState(newstate)
	}
	err = smb.Init()
	if err != nil {
		return
	}
	return nil
}

func NewDefaultStateMachineBody(depth int) StateMachineBody {
	return StateMachineBody{
		States:        map[string]any{},
		_States:       map[string]State{},
		_Depth:        depth,
		_NewStateFunc: NewStateFromMap,
	}
}

func (s *StateMachineBody) SetStartAt(name string) {
	s.StartAt = name
}

func (s *StateMachineBody) GetBone() StateMachineBone {
	bone := StateMachineBone{
		StartAt: s.StartAt,
		States:  make(map[string]StateBone),
	}
	for name, state := range s._States {
		bone.States[name] = state.GetBone()
	}
	return bone
}

// Validate ...
func (s *StateMachineBody) Validate() error {

	// verify StartAt
	if _, ok := s.States[s.StartAt]; !ok {
		return fmt.Errorf("field '%s' state '%s' not found in statemachine",
			StateMachineFieldNames.StartAt, s.StartAt)
	}
	// verify all nodes next in state valid
	for statename, statebone := range s._States {
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
	for _, state := range s._States {
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

func (smb *StateMachineBody) Init() (err error) {

	// verify statemachine
	// 验证startat 在map中

	var stateBonesMap = map[string]int{}
	for _, s := range smb._States {
		stateBonesMap[s.GetName()] = 0
	}

	if _, ok := stateBonesMap[smb.StartAt]; !ok {
		err = fmt.Errorf("field 'StartAt' State '%s' not  found", smb.StartAt)
		return
	}
	// 验证所有节点的next 在state 中
	for _, state := range smb._States {
		for _, next := range state.GetBone().Next {
			if _, ok := stateBonesMap[next]; !ok {
				return fmt.Errorf("state [ %s ] Next '%s' not found in statemachine", state.GetName(), next)
			}
		}
	}

	end := false
	for _, state := range smb._States {
		if state.IsEnd() {
			end = true
			break
		}
	}
	if !end {
		err := fmt.Errorf("statemachine need end state ")
		return err
	}
	return nil
}

// AddState 添加状态机中的状态
func (smb *StateMachineBody) AddState(newstate State) {
	name := newstate.GetName()
	smb._States[name] = newstate
}

// GetGroupStates 生成状态机中所有的state的group信息，并且拉平
func (smb *StateMachineBody) GetGroupStates(startGroupID int) ([]StateMachineGroupState, error) {
	var err error
	var smgss = []StateMachineGroupState{}
	// global group id 全局的group id
	var GlobalGroupID = startGroupID

	// type Process func(sm *StateMachine, groupID int, depth int) error
	type Process func(sm *StateMachineBody, groupId int) error
	var process Process

	process = func(smb *StateMachineBody, groupID int) error {

		var idx = 0
		for name, state := range smb._States {
			idx = idx + 1
			var subgroup *SubGroupState
			newsms := StateMachineGroupState{
				GroupID:    groupID,
				GroupIndex: idx,
				Depth:      smb._Depth,
				Name:       name,
				State:      state,
				SubGroup:   subgroup,
			}
			smgss = append(smgss, newsms)
		}
		return nil
	}

	firstSmb := smb
	err = process(firstSmb, GlobalGroupID)
	return smgss, err
}
