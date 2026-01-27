package dispatcher

import (
	"testing"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backend/config"
	"github.com/skyflow-workflow/skyflow_backend/workflow"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"
	"gopkg.in/go-playground/assert.v1"
)

func TestSchedular_TestRunChoiceEvent(t *testing.T) {

	dbClient, err := rdb.NewDBClient(&rdb.Config{
		DSN:          "mysql://root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
		MaxOpenConns: 100,
		MaxIdleConns: 100,
		MaxIdleTime:  "10s",
		LogLevel:     "info",
	})
	assert.Equal(t, err, nil)

	mockconfig := &config.SkyflowConfig{
		Exporter: &config.ExporterConfig{
			CacheSizeMB: 100,
			PoolSize:    100,
		},
		Dispatcher: &config.DispatcherConfig{
			MaxConcurrency: 100,
			MaxQueueSize:   100,
			Debug:          true,
		},
	}

	workflowSvc, err := workflow.NewWorkflowService(dbClient, queue.NewMockInnerQueue(), mockconfig)
	assert.Equal(t, err, nil)

	dispatcher, err := NewDispatcher(workflowSvc, mockconfig.Dispatcher)
	assert.Equal(t, err, nil)

	msgbody := queue.InnerMessageBody{
		ExecutionID: 5,
		StepID:      15,
		Class:       "Step",
		Type:        "StateNewTurn",
		Data:        "null",
	}
	// msg := queue.NewMockMessage("1234567890", msgbody)

	err = dispatcher.processInnerMessage(msgbody)
	assert.Equal(t, err, nil)

}
