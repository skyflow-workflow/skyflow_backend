package executor

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
)

func TestGetBone(t *testing.T) {

	var myExecutor = StandardExecutor
	var testcases = []struct {
		dbStep po.Step
		except StepBone
		input  string
	}{
		{
			dbStep: po.Step{
				ID:     1,
				Type:   "Task",
				Status: "Running",
				Definition: `{
					"Type": "Task",
					"Resource": "activity:function1",
					"Next": "TMPAlarmShield",
					"Comment": "test comment"
				}`,
				GroupIndex: 1,
				GroupID:    1,
			},
			except: StepBone{
				BaseBone: BaseBone{
					BaseBone: states.BaseBone{
						Type:    "Task",
						Name:    "",
						Next:    []string{"TMPAlarmShield"},
						Comment: "test comment",
					},
					Status: "Running",
					StepID: 1,
					Index:  0,
				},
				StateMachineBone: nil,
				Branches:         nil,
			},
		},
	}

	for _, tt := range testcases {

		step, err := NewExecutionStep(&tt.dbStep, myExecutor)
		assert.Equal(t, err, nil)
		bone := step.GetBone()
		assert.Equal(t, bone, tt.except)
	}

}
