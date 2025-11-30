package executor

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"
)

func TestProcessSuspend(t *testing.T) {
	var testcases = []struct {
		id        int
		wantError bool
	}{
		{
			id:        90,
			wantError: false,
		},
	}

	for _, tt := range testcases {
		state, err := NewSuspendFromID(tt.id, myExecutionService.StandardExecutor)
		fmt.Println(err)
		assert.Equal(t, err != nil, tt.wantError)
		err = state.Run(queue.InnerMessageBody{})
		fmt.Println(err)
		assert.Equal(t, err != nil, tt.wantError)
	}

}
