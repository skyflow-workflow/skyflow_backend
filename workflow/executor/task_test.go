package executor

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"

	"gopkg.in/go-playground/assert.v1"
)

func TestParseTask(t *testing.T) {

	myExecutor := &StandardExecutor
	dbStep := &po.Step{
		Type: "Task",
		Definition: `{
			"Type": "Task",
			"Resource": "abc",
			"Next": "TMPAlarmShield"
		}`,
	}

	exestate, err := NewTaskFromData(dbStep, myExecutor)
	fmt.Println(err)
	fmt.Println(exestate)
}

func TestTask(t *testing.T) {

	myExecutor := &StandardExecutor
	step_id := 52

	taskStep, err := NewTaskFromID(step_id, myExecutor)
	if err != nil {
		fmt.Println(err)
	}
	input, _ := taskStep.GetInput()
	fmt.Println(input)
	fmt.Println(taskStep.State)
	fmt.Println(taskStep.TaskState.Parameters)
	fmt.Println(taskStep.Data)
	inputDatabyte, err := json.Marshal(input)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(inputDatabyte))
	fmt.Println(taskStep.GetBone())
}

func TestRunTask(t *testing.T) {

	var testcases = []struct {
		id        int
		wantError bool
	}{
		{
			id:        147,
			wantError: false,
		},
	}

	for _, tt := range testcases {
		state, err := NewTaskFromID(tt.id, myExecutor)
		fmt.Println(err)
		assert.Equal(t, err != nil, tt.wantError)
		err = state.Run(queue.InnerMessageBody{})
		fmt.Println(err)
		assert.Equal(t, err != nil, tt.wantError)
	}
}

func TestTaskToken(t *testing.T) {
	myExecutor := &StandardExecutor
	token := "e117a41b-f2a4-4195-a259-3a2718327d48"
	task, err := myExecutor.NewTaskFromToken(token, []string{}, nil)
	fmt.Println(err)
	assert.Equal(t, err, nil)
	fmt.Println(task)
	// req := RequestSendTaskFailure{
	// 	Token: ,
	// }
	// task.SendTaskFailure()

}

// func TestSendTaskFailure(t *testing.T) {

// 	myExecutor := &StandardExecutor
// 	// tasktoken := "c3bd022c-1ff0-439a-a7cc-d18194e59818"
// 	tasktoken := "12334w3545mocktoken"
// 	voreq := vo.SendTaskFailureRequest{
// 		TaskToken: tasktoken,
// 		Error:     "test",
// 		Cause:     "test",
// 	}
// 	err := myExecutor.SendTaskFailure(context.Background(), voreq)
// 	fmt.Println(err)

// }
