package executor

import (
	"fmt"
	"testing"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"gopkg.in/go-playground/assert.v1"
)

func TestRunWait(t *testing.T) {

	var testcases = []struct {
		id        int
		wantError bool
	}{
		{
			id:        1605,
			wantError: false,
		},
		{
			id:        1466,
			wantError: false,
		},
	}

	for _, tt := range testcases {
		state, err := NewWaitFromID(tt.id, myExecutionService.StandardExecutor)
		fmt.Println(err)
		assert.Equal(t, err != nil, tt.wantError)
		msg := queue.InnerMessageBody{
			Data: `{"counter":1,"token":""}`,
		}
		err = state.WaitStateWakeup(msg)
		fmt.Println(err)
		assert.Equal(t, err != nil, tt.wantError)
	}
}
