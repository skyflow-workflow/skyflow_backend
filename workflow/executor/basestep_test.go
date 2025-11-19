package executor

import (
	"testing"

	"github.com/skyflow-workflow/skyflow_backend/workflow/po"

	"github.com/go-playground/assert/v2"
)

func TestInitStep(t *testing.T) {
	var err error
	myExecutor := StandardExecutor
	dbStep := &po.Step{
		ID: 1,
		Definition: `{
			"Type": "Task",
			"Name": "Task",
			"Resource": "activity:test_add",
			"Parameters": {
				"a.$" : "$.x",
				"b.$" : "$.y"
			},
			"Next": "NextStep"
		}`,
		Input: `{"x":1,"y":2}`,
	}
	exceptInput := map[string]any{
		"a": float64(1),
		"b": float64(2),
	}
	exeStep, err := NewExecutionStep(dbStep, myExecutor)
	assert.Equal(t, err, nil)
	getInputData, err := exeStep.GetInput()
	assert.Equal(t, err, nil)
	assert.Equal(t, getInputData, exceptInput)
}
