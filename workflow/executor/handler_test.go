package executor

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/skyflow-workflow/skyflow_backbend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
	"gopkg.in/go-playground/assert.v1"
)

func TestStartExecutionStateMachine(t *testing.T) {

	var testcases = []struct {
		definition string
		input      string
	}{{
		definition: `{
			"Comment": "An example of the Amazon States Language using a choice state.",
			"Type": "statemachine",
			"StartAt": "FirstState",
			"TimeoutSeconds": 0,
			"DefaultInput":{
				"x":3,
				"z":10
			},
			"States": {
			  "FirstState": {
				"Type": "Task",
				"Resource": "arn:aws:lambda:REGION:ACCOUNT_ID:function:FUNCTION_NAME",
				"End": true
			  }
			}
		}
`,
		input: `{"x":1, "y":2}`,
	}, {
		definition: `
Version : v1
Type: "pipeline"
# stages 所有阶段
stages:
  - name: 任务准备阶段
    # stage的所有job
    jobs:
      - name: 任务准备
        steps:
          - name: 发送通知消息给任务发起人
            type: "Task"
            resources: "seed/send_im_message"
            parameter:
              receivers : "$.actor"
              chat_id: "$.chat_id"
              content: |
                准备发起seed变更任务,请悉知, {{$.order.url}}`,
		input: `{"x":1, "y":2}`,
	}}

	for idx, tt := range testcases {
		fmt.Printf("process index --- %d \n", idx)

		var req = vo.StartExecutionRequest{
			Title:                  fmt.Sprintf("unittest for create execution : %s", time.Now().String()),
			StateMachineDefinition: tt.definition,
			Input:                  tt.input,
		}
		result, err := myExecutionService.StartExecution(req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(result)

	}
}
func TestStartExecutionPipeline(t *testing.T) {

	var piplepath = "../../examples/pipeline/jobmap_pipeline.yaml"
	var input = `{"x":1, "y":2}`

	// var piplepath = "../../examples/pipeline/task_pipeline.yaml"
	// var input = `{"x":1, "y":2}`

	content, err := os.ReadFile(piplepath)
	if err != nil {
		log.Println(err)
		return
	}
	var req = vo.StartExecutionRequest{
		Title:                  fmt.Sprintf("create for create execution : %s", time.Now().String()),
		StateMachineDefinition: string(content),
		Input:                  input,
	}
	result, err := myExecutionService.StartExecution(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)

}

func TestSendTaskSkip(t *testing.T) {

	step_id := 2479
	task, err := NewTaskFromID(step_id, myExecutionService.StandardExecutor)
	fmt.Println(err)
	assert.Equal(t, err, nil)
	fmt.Println(task)

	ctx := context.Background()
	req := vo.SendStepSkipRequest{
		StepID:       step_id,
		NextStepName: "MHello",
		Output:       `{"xyz":"xx"}`,
	}
	err = myExecutionService.SendStepSkip(ctx, req)
	fmt.Println(err)
	assert.Equal(t, err, nil)

	// req := RequestSendTaskFailure{
	// 	Token: ,
	// }
	// task.SendTaskFailure()

}

func TestCreateUUID(t *testing.T) {

	uuidstr, err := toolkit.CreateUUID()
	fmt.Println(err)
	assert.Equal(t, err == nil, true)
	fmt.Println(uuidstr)

}

func TestGetActivityTask(t *testing.T) {

	ctx := context.Background()
	for i := 0; i < 3; i++ {

		uri := "activity:unittest/add"
		time.Sleep(1 * time.Second)
		resp, err := myExecutionService.GetActivityTask(ctx, vo.GetActivityTaskRequest{
			ActivityURI: uri,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(resp)
	}

}

func TestSendTaskSuccess(t *testing.T) {
	token := "e117a41b-f2a4-4195-a259-3a2718327d48"
	output := "3"

	ctx := context.Background()

	err := myExecutionService.SendTaskSuccess(ctx, vo.SendTaskSuccessRequest{
		TaskToken: token,
		Output:    output,
	})
	if err != nil {
		t.Error("err: ", err)
		return
	}
}

func TestChangeStepGroupStatus(t *testing.T) {

	stepgroup_id := 294
	err := myExecutionService.ChangeStepGroupStatus(stepgroup_id, ExecutionStatus.Running, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
func TestResumeStep(t *testing.T) {

	step_id := 303
	err := myExecutionService.ResumeSuspendingStep(step_id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func TestResumeExecution(t *testing.T) {

	exeid := 16
	err := myExecutionService.ResumeExecution(exeid)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func TestSendStepRetry(t *testing.T) {

	var step_id = 749
	err := myExecutionService.SendStepRetry(step_id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
