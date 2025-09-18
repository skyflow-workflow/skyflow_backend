package states

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseStateGroup(t *testing.T) {

	var testcases = []struct {
		definition string
	}{
		{
			definition: `
			{
				"StartAt": "LookupAddress",
				"States": {
				  "LookupAddress": {
					"Type": "Task",
					"Resource":
					  "arn:aws-cn:lambda:us-east-1:123456789012:function:AddressFinder",
					"End": true
				  }
				}
			}
			`,
		}, {
			definition: `
			{
				"StartAt": "LookupPhone",
				"States": {
				  "LookupPhone": {
					"Type": "Task",
					"Resource":
					  "arn:aws-cn:lambda:us-east-1:123456789012:function:PhoneFinder",
					"End": true
				  }
				}
			  }`,
		},
	}

	for idx, testcase := range testcases {
		fmt.Println("idx -- ", idx)
		pl, err := NewStateGroupFromString(testcase.definition, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(pl)
		bone := pl.GetBone()
		bonebyte, err := json.Marshal(bone)
		fmt.Println(string(bonebyte))
	}
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
