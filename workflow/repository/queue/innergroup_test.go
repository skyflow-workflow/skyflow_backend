package queue

import (
	"log/slog"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/skyflow-workflow/skyflow_backbend/mock"
)

func TestCreateInnerQueueGroup(t *testing.T) {

	masterQueue, err := NewInnerMessageQueueFromConfig(mock.LocalUnitTestKafkaDSN)
	assert.Equal(t, nil, err)

	delayQueue := NewDBMessageQueue(mock.MockDBClient, DefaultDBDelayQueueOption)

	queuegroup, err := NewInnerQueueGroup(masterQueue, delayQueue)
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
