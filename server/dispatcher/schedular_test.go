package dispatcher

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mmtbak/microlibrary/config"
	"github.com/mmtbak/microlibrary/mq"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
	"trpc.group/trpc-go/tnet/log"
)

var dbrepo *rdb.DBRepository
var myExecutionService execution.ExecutionService

func init() {
	conf := config.AccessPoint{
		Source:  "memory:///?size=10",
		Options: map[string]interface{}{},
	}

	basequeue, err := mq.NewMessageQueue(conf)
	if err != nil {
		log.Error(err.Error())
		return
	}
	innerqueue := queue.NewInnerMessageQueue2(basequeue, queue.Option{Debug: true, Logger: slog.Default()})

	myExecutionService = execution.NewExecutionService(dbrepo, innerqueue)

}
func TestExecutionWorker(t *testing.T) {

	// ExecutionWorkerManger()
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func TestExecuteTask(t *testing.T) {

	var testcases = []vo.StartExecutionRequest{
		{
			WorkflowDefinition: `
			{
				"StartAt":"P1S",
				"States":{
					"P1S":{
						"Type":"Pass",
						"End":true
					}
				}
			}`,
			Input: `{}`,
		},
		{
			WorkflowDefinition: `
			{
				"StartAt":"P1S",
				"States":{
					"P1S":{
						"Type":"Task",
						"Resource":"activity:om2/function",
						"End":true
					}
				}
			}`,
			Input: `{}`,
		},
	}

	for _, req := range testcases {
		dbexe, err := myExecutionService.StartExecution(req)
		fmt.Println(err)
		fmt.Println(dbexe)
		assert.Equal(t, err == nil, true)
	}
	select {}
}
