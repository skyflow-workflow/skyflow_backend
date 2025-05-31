package queue

import (
	"log/slog"
	"testing"

	"github.com/skyflow-workflow/skyflow_backbend/mock"
)

func TestCreateInnerQueueGroup(t *testing.T) {

	queuegroup, err := NewInnerQueueGroupFromConfig(mock.LocalUnitTestKafka, mock.LocalUnitTestMySQLConfig)
	if err != nil {
		t.Error(err)
		return
	}

	var testcases = []struct {
		Message InnerMessageBody
	}{
		{
			Message: InnerMessageBody{
				ExecutionID: 1,
				Class:       "Execution",
				Type:        "ExecutionInit",
				StepID:      1,
				Data:        `{"testkey":"testvalue"}`,
			},
		},
		{
			Message: InnerMessageBody{
				ExecutionID: 1,
				Class:       "Execution",
				Type:        "ExecutionInit",
				StepID:      2,
				Data:        `{"testkey":"testvalue"}`,
			},
		},
		{
			Message: InnerMessageBody{
				ExecutionID: 2,
				Class:       "Execution",
				Type:        "ExecutionInit",
				StepID:      2,
				Data:        `{"testkey":"testvalue"}`,
			},
		},
		{
			Message: InnerMessageBody{
				ExecutionID: 2,
				Class:       "Execution",
				Type:        "ExecutionInit",
				StepID:      3,
				Data:        `{"testkey":"testvalue"}`,
			},
		},
		{
			Message: InnerMessageBody{
				ExecutionID: 2,
				Class:       "Step",
				Type:        "StepInit",
				StepID:      4,
				Data:        `{"testkey":"testvalue"}`,
			},
		},
		{
			Message: InnerMessageBody{
				ExecutionID: 2,
				Class:       "Execution",
				Type:        "ExecutionInit",
				StepID:      5,
				Data:        `{"testkey":"testvalue"}`,
			},
		},
	}

	for idx, tt := range testcases {

		slog.Info("send message idx - ", "index", idx)
		err = queuegroup.SendInnerMessage(tt.Message, nil)
		if err != nil {
			t.Error(err)
			return
		}
	}

}
