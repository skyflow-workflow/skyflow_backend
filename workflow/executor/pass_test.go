package executor

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"
)

func TestExecutePass(t *testing.T) {

	var testcases = []struct {
		id        int
		wantError bool
	}{
		{
			id:        1650,
			wantError: false,
		},
	}

	for _, tt := range testcases {
		state, err := NewPassFromID(tt.id, myExecutionService.StandardExecutor)
		fmt.Println(err)
		assert.Equal(t, err != nil, tt.wantError)
		err = state.Run(queue.InnerMessageBody{})
		fmt.Println(err)
		assert.Equal(t, err != nil, tt.wantError)
	}

}

func TestPass_ExecutePass(t *testing.T) {
	var dbStep = po.Step{
		ID:           1,
		ExecutionID:  17,
		ExecuteIndex: 0,
		GroupID:      1,
		Name:         "Hello",
		GroupIndex:   1,
		Status:       "WaitInit",
		Depth:        1,
		Type:         "Pass",
		Resource:     "",
		Definition:   "\"{\"Name\":\"Hello\",\"Type\":\"Pass\",\"OutputPath\":\"$\",\"MaxExecuteTimes\":1000,\"End\":true,\"Result\":{\"message\":\"Hello World\"}}\"",
		Input:        "{}",
		Output:       "",
		Exception:    "",
		Data:         "{}",
		StartTime:    nil,
		FinishTime:   nil,
		ExpireTime:   nil,
		ExecuteCount: 0,
		UpdateTime:   time.Now(),
		CreateTime:   time.Now(),
	}

	state, err := NewPassFromData(&dbStep, myExecutionService.StandardExecutor)
	fmt.Println(err)
	assert.Equal(t, err != nil, false)
	err = state.Run(queue.InnerMessageBody{})
	fmt.Println(err)
	assert.Equal(t, err != nil, false)

}
