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
	StartAt string
	States  map[string]State
}

// NewStateMachineBodyFromMap Parse data to StateMachine
func NewStateMachineBodyFromMap(data map[string]interface{}, depth int) (smb StateMachineBody, err error) {

	// 状态机深度不能超过最大深度
	if depth > MaxDepth {
		err = fmt.Errorf("statemachine depth  reach MaxDepth [ %d ] ", MaxDepth)
		return
	}
	smb = NewStateMachineBody(depth)

	err = smb.InitByMap(data)
	if err != nil {
		return
	}
	return
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

// InitByMap Use map to init StateMachineBody content
func (smb *StateMachineBody) InitByMap(data map[string]interface{}) (err error) {

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
		content, ok := stateObj.(map[string]interface{})
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
func (smb *StateMachineBody) Init() (err error) {

	// verify statemachine
	// 验证startat 在map中

	var statebonesmap = map[string]int{}
	for _, s := range smb.States {
		statebonesmap[s.GetName()] = 0
	}

	if _, ok := statebonesmap[smb.StartAt]; !ok {
		err = fmt.Errorf("field 'StartAt' State '%s' not  found", smb.StartAt)
		return
	}
	// 验证所有节点的next 在state 中
	for _, state := range smb.States {
		for _, next := range state.GetBone().Next {
			if _, ok := statebonesmap[next]; !ok {
				return fmt.Errorf("state [ %s ] Next '%s' not found in statemachine", state.GetName(), next)
			}
		}
	}

	end := false
	for _, state := range smb.States {
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
	smb.States[name] = newstate
}

// GetGroupStates 生成状态机中所有的state的group信息，并且拉平
func (smb *StateMachineBody) GetGroupStates() ([]StateMachineGroupState, error) {
	var err error
	var smgss = []StateMachineGroupState{}
	// var startdepth = StartDepth
	var GGroupID = StartGroupID

	// type Process func(sm *StateMachine, groupID int, depth int) error
	type Process func(sm *StateGroup, groupID int) error
	var process Process

	process = func(sm *StateGroup, groupID int) error {

		for idx, state := range sm.States {
			name := state.GetName()
			var subgroup *SubGroupState
			stype := state.GetType()
			if stype == string(StateTypes.Parallel) {
				ps, ok := state.(*ParallelState)
				if !ok {
					err = fmt.Errorf("transform state '%s' to  'Parallel' failed", name)
					return err
				}
				// 深度+1, deprecated, 深度计算在初始化的时候就算了
				// var newdepth = depth + 1
				//
				GGroupID = GGroupID + 1
				//SubGroupID  子组的GroupID
				SubGroupID := GGroupID

				for idx2, branchsm := range ps._Branches {
					// 每次都是新的group
					// newGrouID 每个组内的group id
					GGroupID = GGroupID + 1
					subgrouID := GGroupID

					err = process(branchsm, subgrouID)
					if err != nil {
						return err
					}

					smstate := StateMachineGroupState{
						GroupID:    SubGroupID,
						GroupIndex: idx2,
						Depth:      branchsm._Depth,
						Name:       branchsm.Name,
						State:      branchsm.BaseState,
						SubGroup: &SubGroupState{
							SubGroupID:    subgrouID,
							SubStartAt:    branchsm.StartAt,
							MasterGroupID: groupID,
							MasterName:    name,
						},
					}
					smgss = append(smgss, smstate)

				}
			}
			newsms := StateMachineGroupState{
				GroupID:    groupID,
				GroupIndex: idx,
				Depth:      sm._Depth,
				Name:       name,
				State:      state,
				SubGroup:   subgroup,
			}
			smgss = append(smgss, newsms)
		}
		return nil
	}
	outersg := &StateGroup{
		StateMachineBody: *smb,
	}
	err = process(outersg, GGroupID)
	return smgss, err
}
