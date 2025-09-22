package states

// StateMachine ...
type StateMachine struct {
	*StateMachineHeader
	*StateMachineBody
}

func NewStateMachineFromString(data string) (*StateMachine, error) {

	mapdata, err := ToMap(data)
	if err != nil {
		return nil, err
	}
	sm, err := NewStateMachineFromMap(mapdata)
	if err != nil {
		return nil, err
	}
	return sm, err
}

func NewStateMachineFromMap(data map[string]interface{}) (*StateMachine, error) {

	var err error
	header, err := NewStateMachineHeaderFromMap(data)
	if err != nil {
		return nil, err
	}
	body, err := NewStateMachineBodyFromMap(data, StartDepth)
	if err != nil {
		return nil, err
	}
	sm := &StateMachine{
		StateMachineHeader: header,
		StateMachineBody:   body,
	}

	return sm, err
}

func (sm *StateMachine) Init() (err error) {
	return nil
}

// GetInput  GetInput
func (sm StateMachine) GetInput(input any, executioninfo any) (any, error) {
	return sm.StateMachineHeader.GetInput(input, executioninfo)
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
	bone := sm.StateMachineBody.GetBone()
	return bone
}

func (sm *StateMachine) GetGroupStates() ([]StateMachineGroupState, error) {
	groupstates, err := sm.StateMachineBody.GetGroupStates(StartGroupID)
	if err != nil {
		return nil, err
	}
	return groupstates, nil
}
