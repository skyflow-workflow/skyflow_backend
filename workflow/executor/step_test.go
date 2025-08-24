package executor

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
)

func TestGetBone(t *testing.T) {

	var myExecutor = StandardExecutor
	var testcases = []struct {
		step_id int
		input   string
	}{
		{
			step_id: 223,
		},
		{
			step_id: 402,
		},
	}

	for _, tt := range testcases {

		fmt.Println(tt.step_id)
		step, err := NewExecutionStep(&po.Step{}, myExecutor)
		assert.Equal(t, err, nil)
		bone := step.GetBone()
		fmt.Println(bone)
	}

}
