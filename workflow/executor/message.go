package executor

import (
	"github.com/skyflow-workflow/skyflow_backbend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
)

// message.go defines message struct processed in workflow processing

// TaskTimeoutMessage task timeout message
type TaskTimeoutMessage struct {
	HeartbeatCount int    `json:"heartbeat_count"`
	ExecuteCount   int    `json:"execute_count"`
	TaskToken      string `json:"task_token"`
}

// TaskWakeupMessage  task wakeup message
// TaskWakeupMessage is used to wake up a task state in retry strategy or retry by user manually
type TaskWakeupMessage struct {
	Counter int    `json:"counter"`
	Token   string `json:"token"`
}

type StepExecuteMessage struct {
	// if block is true, the step will be blocked
	// default is false
	Block bool `json:"block"`
}

// NewExecutionMessage create a new message
func NewExecutionMessage(execution_id int, types string, data interface{}) queue.InnerMessageBody {

	s, _ := toolkit.EncodeToString(data)
	message := queue.InnerMessageBody{
		ExecutionID: execution_id,
		Type:        types,
		StepID:      0,
		Data:        s,
		Class:       queue.MessageClass.Execution,
	}
	return message
}

// NewStepMessage create a new step message
func NewStepMessage(execution_id int, types string, step_id int, data interface{}) queue.InnerMessageBody {

	s, _ := toolkit.EncodeToString(data)
	message := queue.InnerMessageBody{
		ExecutionID: execution_id,
		Type:        types,
		StepID:      step_id,
		Data:        s,
		Class:       queue.MessageClass.Step,
	}
	return message
}
