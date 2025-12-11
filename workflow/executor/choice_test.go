package executor

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"
)

func TestRunChoice(t *testing.T) {

	var testcases = []struct {
		id        int
		wantError bool
	}{
		{
			id:        1821,
			wantError: false,
		},
	}

	for _, tt := range testcases {
		state, err := NewChoiceFromID(tt.id, myExecutionService.StandardExecutor)
		fmt.Println(err)
		assert.Equal(t, err != nil, tt.wantError)
		err = state.Run(queue.InnerMessageBody{})
		fmt.Println(err)
		assert.Equal(t, err != nil, tt.wantError)
	}

}

func TestChoice_RunNumberEqueal(t *testing.T) {

	dbStep := &po.Step{
		ID: 1821,
		Definition: `{
			"Type": "Choice",
			"Choices": [
			  {
				"Variable": "$.value",
				"NumericEquals": 0,
				"Next": "ValueIsZero"
			  }
			],
			"Default": "DefaultState"
		  }`,
		Input: `{"value":200,"b":100,"z":300}`,
	}

	step, err := NewChoiceFromData(dbStep, myExecutionService.StandardExecutor)
	assert.Equal(t, err != nil, false)
	err = step.Run(queue.InnerMessageBody{})
	assert.Equal(t, err != nil, false)

}
