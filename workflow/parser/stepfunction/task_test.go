package stepfunction

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
)

func TestDecodeTask(t *testing.T) {
	var testcases = []struct {
		definition string
		expected   *states.Task
	}{
		{
			definition: `{
				"Type":"Task",
				"Comment": "add task",
				"Resource": "activity:unittest/add",
				"Next": "",
				"End": false
			}
			`,
			expected: &states.Task{
				BaseState: &states.BaseState{
					Type:            "Task",
					Comment:         "add task",
					MaxExecuteTimes: 1000,
					InputPath:       "$",
					OutputPath:      "$",
					ResultPath:      "",
					End:             false,
					Next:            "",
				},
				TaskBody: &states.TaskBody{
					Resource:         "activity:unittest/add",
					HeartbeatSeconds: 0,
					TimeoutSeconds:   0,
					Catch:            []states.CatchNode{},
					Retry:            []states.RetryNode{},
				},
			},
		},
		{
			definition: `{
				"Type": "Task",
				"Parameters": {
					"x.$": "$.x",
					"y.$": "$.y"
				},
				"Resource": "activity:unittest/Testadd",
				"ResultPath": "$.z",
				"Next": "S2"
			}
			`,
			expected: &states.Task{
				BaseState: &states.BaseState{
					Type:    "Task",
					Comment: "",
					Parameters: map[string]interface{}{
						"x.$": "$.x",
						"y.$": "$.y",
					},
					MaxExecuteTimes: 1000,
					InputPath:       "$",
					OutputPath:      "$",
					ResultPath:      "$.z",
					End:             false,
					Next:            "S2",
				},
				TaskBody: &states.TaskBody{
					Resource:         "activity:unittest/Testadd",
					HeartbeatSeconds: 0,
					TimeoutSeconds:   0,
					Catch:            []states.CatchNode{},
					Retry:            []states.RetryNode{},
				},
			},
		},
	}

	decoder := NewStepfuncionDecoder(nil, nil)

	for i, tt := range testcases {
		t.Run(fmt.Sprintf("index_%d", i), func(t *testing.T) {
			var err error
			err = tt.expected.BaseState.Init()
			assert.Equal(t, err, nil)
			err = tt.expected.TaskBody.Init()
			assert.Equal(t, err, nil)

			err = tt.expected.Init()
			assert.Equal(t, err, nil)
			state, err := decoder.DecodeStateDefintion(tt.definition)
			assert.Equal(t, err, nil)
			opts := cmp.Options{cmpopts.IgnoreUnexported(states.BaseState{}), cmpopts.IgnoreUnexported(states.TaskBody{})}
			diff := cmp.Diff(state, tt.expected, opts...)
			assert.Equal(t, "", diff)

		})
	}
}
