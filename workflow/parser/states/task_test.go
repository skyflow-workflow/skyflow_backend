package states

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestTaskGetBone(t *testing.T) {

	var testcases = []struct {
		task    *TaskState
		except  StateBone
		wantErr bool
	}{
		{
			task: &TaskState{
				BaseState: &BaseState{
					Type: "Task",
					Next: "NextState",
				},
				TaskBody: &TaskBody{
					Catch: []TaskCatchNode{
						{
							Next: "CatchNextState",
						},
						{
							Next: "CatchNextState",
						},
					},
				},
			},
			except: StateBone{
				BaseBone: BaseBone{
					Type: "Task",
					Next: []string{"NextState", "CatchNextState", "CatchNextState"},
				},
			},
			wantErr: false,
		},
		{
			task: &TaskState{
				BaseState: &BaseState{
					Type: "Task",
					Next: "NextState",
				},
				TaskBody: &TaskBody{
					Catch: []TaskCatchNode{
						{
							Next: "CatchNextState",
						},
						{
							Next: "CatchNextState",
						},
					},
				},
			},
			except: StateBone{
				BaseBone: BaseBone{
					Type: "Task",
					Next: []string{"NextState", "CatchNextState", "CatchNextState"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range testcases {
		actualbone := tt.task.GetBone()
		assert.Equal(t, actualbone, tt.except)
	}

}

func TestTaskGetTaskTimeout(t *testing.T) {

}

func TestTaskGetOutput(t *testing.T) {

}

func TestTaskHasIntersection(t *testing.T) {

}

func TestTaskGetOutputWithPath(t *testing.T) {

}

func TestHasIntersection(t *testing.T) {

	var testcases = []struct {
		a      []string
		b      []string
		except bool
	}{
		{
			a:      []string{"a", "b", "c"},
			b:      []string{"a", "b", "c"},
			except: true,
		},
		{
			a:      []string{"a", "b", "c"},
			b:      []string{"a", "b", "c", "d"},
			except: true,
		},
		{
			a:      []string{"a", "b", "c"},
			b:      []string{"a", "b", "c", "d", "e"},
			except: true,
		},
		{
			a:      []string{"a", "b", "c"},
			b:      []string{"d", "e", "f"},
			except: false,
		},
	}

	for _, tt := range testcases {
		actual := HasIntersection(tt.a, tt.b)
		assert.Equal(t, actual, tt.except)
	}

}

func TestParseTaskState(t *testing.T) {

	var testcases = []struct {
		template  string
		wantError bool
	}{
		{
			template: `
	{
		"Type": "Task",
		"Resource": "activity:nest/polaris_del_instance_by_ip",
		"Parameters": {
		  "ip.$": "$.ip"
		},
		"Retry": [
		  {
			"ErrorEquals": [
			  "DeleteInstanceFaild"
			],
			"IntervalSeconds": 600,
			"MaxAttempts": 10,
			"BackoffRate": 1
		  },
		  {
			"ErrorEquals": [
			  "DeleteInstanceFaild"
			],
			"IntervalSeconds": 600,
			"MaxAttempts": 10
		  }
		],
		"Catch":[
			{
				"ErrorEquals": [
				  "DeleteInstanceFaild"
				],
				"Next": "A"
			},
			{
				"ErrorEquals": [
				  "DeleteInstanceFaild"
				],
				"Next": "B"
			}
		],
		"ResultPath": "",
		"Comment": "删除业务在北极星上的注册",
		"Next": "TMPAlarmShield"
	}`,
			wantError: false,
		},
		{
			template: `	{
				"Type": "Task",
				"Resource": "",
				"Next": "TMPAlarmShield"
			}
			`,
			wantError: true,
		},
		{
			template: `	{
				"Type": "Task",
				"Resource": "abc",
				"Next": "TMPAlarmShield",
				"Bone": "xxx"
			}
			`,
			wantError: false,
		},
		{
			template: `	{
				"Type": "Task",
				"Resource": "abc",
				"Next": "TMPAlarmShield",
				"InputPath": 1233
			}
			`,
			wantError: true,
		},
		{
			template: `	{
				"Type": "Task",
				"Resource": "abc",
				"Next": "xyz",
				"_InputPath": 1233
			}
			`,
			wantError: false,
		},
		{
			template: `	{
				"Type": "Task",
				"Resource": "abc",
				"Next": "xyz",
				"TimeoutSeconds": 1233
			}
			`,
			wantError: false,
		},
		{
			template: `	{
				"Type": "Task",
				"Resource": "abc",
				"Next": "xyz",
				"TimeoutSeconds": -1
			}
			`,
			wantError: true,
		},
		{
			template: `	{
				"Type": "Task",
				"Resource": "arn:aws:lambda:REGION:ACCOUNT_ID:function:FUNCTION_NAME",
				"Retry2":-11.2,
				"Parameters" : {
					"abc.$" : "$.foo",
					"123": [{
						"var1.$" : "$.boo",
						"var2.$" : "$.foo"
					}
					]
				},
				"Retry":[
					{
					  "ErrorEquals": [ "ErrorA", "ErrorB" ],
					  "IntervalSeconds": 1,
					  "BackoffRate": 2.0,
					  "MaxAttempts": 2
					},
					{
					  "ErrorEquals": [ "ErrorC" ],
					  "IntervalSeconds": 5
					}
				  ],
				  "Catch": [
					{
					  "ErrorEquals": [ "States.ALL" ],
					  "Next": "Z"
					}
				  ],
				"Next": "ChoiceState"
			  }`,
			wantError: false,
		},
	}

	for _, tt := range testcases {

		var mapdata map[string]interface{}
		err := json.Unmarshal([]byte(tt.template), &mapdata)
		if err != nil {
			fmt.Println(err)
			continue
		}
		state, err := NewTaskStateFromMap(mapdata)
		assert.Equal(t, tt.wantError, err != nil)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(state)

	}

}

func TestParseTask(t *testing.T) {
	taskdefintion := `{
		"Type": "Task",
		"Resource": "arn:aws:lambda:REGION:ACCOUNT_ID:function:FUNCTION_NAME",
		"InputPath" : "$.mydata",
		"Parameters" : {
			"abc.$" : "$.foo",
			"123": [{
				"var1.$" : "$.boo",
				"var2.$" : "$.foo"
			}
			]
		},
		"MaxExecuteTimes":10,
		"Retry":[
			{
			  "ErrorEquals": [ "ErrorA", "ErrorB" ],
			  "IntervalSeconds": 1,
			  "BackoffRate": 2.0,
			  "MaxAttempts": 2
			},
			{
			  "ErrorEquals": [ "ErrorC" ],
			  "IntervalSeconds": 5
			}
		  ],
		  "Catch": [
			{
			  "ErrorEquals": [ "States.ALL" ],
			  "Next": "Z"
			}
		  ],
		"Next": "ChoiceState"
	  }`

	taskstate, err := NewTaskStateFromString(taskdefintion)
	assert.Equal(t, err, nil)
	assert.Equal(t, taskstate.Resource, "arn:aws:lambda:REGION:ACCOUNT_ID:function:FUNCTION_NAME")
	assert.Equal(t, taskstate.InputPath, "$.mydata")
	assert.Equal(t, taskstate.Parameters, map[string]interface{}{
		"abc.$": "$.foo",
		"123": []interface{}{
			map[string]interface{}{
				"var1.$": "$.boo",
				"var2.$": "$.foo",
			},
		},
	})
	assert.Equal(t, taskstate.MaxExecuteTimes, 10)
	assert.Equal(t, taskstate.TaskBody.Retry, []TaskRetryNode{
		{
			ErrorEquals:     []string{"ErrorA", "ErrorB"},
			IntervalSeconds: 1,
			BackoffRate:     2.0,
			MaxAttempts:     2,
		},
		{
			ErrorEquals:     []string{"ErrorC"},
			IntervalSeconds: 5,
			BackoffRate:     1.5,
			MaxAttempts:     3,
		},
	})
	assert.Equal(t, taskstate.TaskBody.Catch, []TaskCatchNode{
		{
			ErrorEquals: []string{"States.ALL"},
			Next:        "Z",
			ResultPath:  "$",
		},
	})
	assert.Equal(t, taskstate.Next, "ChoiceState")
	assert.Equal(t, taskstate.GetBone(), StateBone{
		BaseBone: BaseBone{
			Type:    "Task",
			Name:    "",
			Next:    []string{"ChoiceState", "Z"},
			End:     false,
			Comment: "",
		},
	})
}

func TestRunTask(t *testing.T) {
	taskjson := `{
		"Type": "Task",
		"Resource": "arn:aws:lambda:REGION:ACCOUNT_ID:function:FUNCTION_NAME",
		"InputPath" : "$.mydata",
		"Parameters" : {
			"abc.$" : "$.foo",
			"123": [{
				"var1.$" : "$.boo",
				"var2.$" : "$.foo"
			}
			]
		},
		"MaxExecuteTimes":10,
		"Retry":[
			{
			  "ErrorEquals": [ "ErrorA", "ErrorB" ],
			  "IntervalSeconds": 1,
			  "BackoffRate": 2.0,
			  "MaxAttempts": 2
			},
			{
			  "ErrorEquals": [ "ErrorC" ],
			  "IntervalSeconds": 5
			}
		  ],
		  "Catch": [
			{
			  "ErrorEquals": [ "States.ALL" ],
			  "Next": "Z"
			}
		  ],
		"Next": "ChoiceState"
	  }`

	stateInput := map[string]interface{}{
		"mydata": map[string]interface{}{
			"foo": 2.0,
			"boo": 2,
		},
	}
	exceptTaskInput := map[string]interface{}{
		"abc": 2.0,
		"123": []interface{}{
			map[string]interface{}{
				"var1": 2,
				"var2": 2.0,
			},
		},
	}

	task, err := NewTaskStateFromString(taskjson)
	assert.Equal(t, err, nil)
	taskInput, err := task.GetInput(stateInput)
	assert.Equal(t, err, nil)
	assert.Equal(t, taskInput, exceptTaskInput)
}
