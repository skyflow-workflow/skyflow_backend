package stepfunction

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
)

func TestDecodeTask(t *testing.T) {
	var testcases = []struct {
		name       string
		definition string
		expected   *states.Task
		wantError  error
	}{
		{
			name: "task",
			definition: `{
				"Type":"Task",
				"Comment": "add task",
				"Resource": "activity:unittest/add",
				"Next": ""
			}
			`,
			expected: &states.Task{
				BaseState: &states.BaseState{
					Type:            "Task",
					Comment:         "add task",
					MaxExecuteTimes: 1000,
					InputPath:       "$",
					OutputPath:      "$",
					ResultPath:      "$",
					End:             false,
					Next:            "",
				},
				TaskBody: &states.TaskBody{
					Resource:         "activity:unittest/add",
					HeartbeatSeconds: 0,
					TimeoutSeconds:   0,
					Catch:            []states.TaskCatchNode{},
					Retry:            []states.TaskRetryNode{},
				},
			},
		},
		{
			name: "task with parameters",
			definition: `{
				"Type": "Task",
				"Parameters": {
					"x.$": "$.x",
					"y.$": "$.y"
				},
				"Resource": "activity:unittest/add",
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
					Resource:         "activity:unittest/add",
					HeartbeatSeconds: 0,
					TimeoutSeconds:   0,
					Catch:            []states.TaskCatchNode{},
					Retry:            []states.TaskRetryNode{},
				},
			},
		},
	}

	decoder := NewStepfuncionDecoder(nil, nil)

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			state, err := decoder.DecodeStateDefintion(tt.definition)
			if err != nil {
				assert.Equal(t, tt.wantError != nil, true)
				assert.Equal(t, err.Error(), tt.wantError.Error())
				return
			}
			opts := cmp.Options{cmpopts.IgnoreUnexported(states.BaseState{}), cmpopts.IgnoreUnexported(states.TaskBody{})}
			diff := cmp.Diff(state, tt.expected, opts...)
			assert.Equal(t, "", diff)

		})
	}
}
