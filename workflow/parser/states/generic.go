package states

import "encoding/json"

// GState  Grammer Gerneric State 单步骤类型泛型
type GState interface {
	TaskState | ChoiceState | WaitState | SucceedState | FailState | SuspendState
}

// GGroupState Grammer Generic Group State 有子流程的类型。
type GGroupState interface {
	MapState | ParallelState | StateGroup
}

// NewGenericStateFromString   Create New Generic Single Step From String
func NewGenericStateFromString[GT GState](definition string, f func(map[string]interface{}) (*GT, error)) (state *GT, err error) {
	var data = map[string]interface{}{}
	err = json.Unmarshal([]byte(definition), &data)
	if err != nil {
		return
	}
	state, err = f(data)
	return
}

// NewGenericStateFromString  Create NewGeneric Group State
func NewGenericGroupStateFromString[GT GGroupState](definition string, depth int, f func(map[string]interface{}, int) (*GT, error)) (state *GT, err error) {
	var data = map[string]interface{}{}
	err = json.Unmarshal([]byte(definition), &data)
	if err != nil {
		return
	}
	state, err = f(data, depth)
	return
}
