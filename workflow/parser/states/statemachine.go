package states

// StateMachine ...
type StateMachine struct {
	*StateMachineHeader
	*StateMachineBody
}

func NewStateMachine() *StateMachine {
	sm := StateMachine{
		StateMachineHeader: NewDefautStateMachineHeader(),
		StateMachineBody:   NewStateMachineBody(StartDepth),
	}
	return &sm
}

// InitByMap Inititalize 初始化一个StateMachine 实例
// jsonMap 一个map[string]interface[] 初始化数据源
// NOCC:golint/fnsize("设计如此")
func (sm *StateMachine) InitByMap(data map[string]interface{}) (err error) {

	defer func() {
		if p := recover(); p != nil {
			err = p.(error)
		}
	}()

	err = sm.StateMachineHeader.InitByMap(data)
	if err != nil {
		return
	}
	err = sm.StateMachineBody.InitByMap(data)
	if err != nil {
		return
	}

	err = sm.Init()
	if err != nil {
		return
	}

	return err
}

func (sm *StateMachine) Init() (err error) {

	return sm.StateMachineBody.Init()
}

// GetInput  GetInput
func (sm StateMachine) GetInput(input interface{}, executioninfo interface{}) (interface{}, error) {
	return sm.StateMachineHeader.GetInput(input, executioninfo)
}

func NewStateMachineBody(depth int) StateMachineBody {
	return StateMachineBody{
		States:        map[string]interface{}{},
		_States:       []State{},
		_Depth:        depth,
		_NewStateFunc: NewStateFromMap,
		_Bone: StateMachineBone{
			StartAt: "",
			States:  map[string]StateBone{},
		},
	}
}

// Validate ...
func (s *StateMachine) Validate() string {
	return "StateMachine"
}

// ParserStateMachine ...
func ParserStateMachine(definition string) error {
	// Parse the state machine
	return nil
}

// GetBone  GetBone
func (sm StateMachine) GetBone() StateMachineBone {
	return sm._Bone
}
