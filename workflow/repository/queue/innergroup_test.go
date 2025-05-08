package queue

import (
	"fmt"
	"testing"

	"github.com/mmtbak/microlibrary/config"
)

func TestCreateInnerQueueGroup(t *testing.T) {

	nornalqueueconfig := config.AccessPoint{
		Source: "kafka://10.22.24.3:9092/?topics=cloudflow_dev&consumergroup=my-event-group&numpartition=1&numreplica=1&version=2.8.1&inital=oldest",
	}

	delayqueueconfig := config.AccessPoint{
		Source: "mysql://neo:Neo!@#123@tcp(10.234.42.154:3306)/cloudflow_dev?charset=utf8mb4&parseTime=True&loc=Local",
		Options: map[string]interface{}{
			"sqllevel": "info",
		},
	}

	queuegroup, err := NewInnerQueueGroupFromConfig(nornalqueueconfig, delayqueueconfig)
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

		fmt.Println("send message idx - ", idx)
		err = queuegroup.SendInnerMessage(tt.Message, nil)
		if err != nil {
			t.Error(err)
			return
		}
	}

}
