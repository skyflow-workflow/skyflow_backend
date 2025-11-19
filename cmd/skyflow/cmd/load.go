package cmd

import (
	"fmt"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backend/config"
	"github.com/skyflow-workflow/skyflow_backend/workflow"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"

	tconfig "trpc.group/trpc-go/trpc-go/config"
)

func LoadConfig(customConfigFilePath string) (*config.SkyflowConfig, error) {

	customCfg, err := tconfig.Load(customConfigFilePath, tconfig.WithCodec("yaml"))
	if err != nil {
		return nil, fmt.Errorf("error loading custom configuration file: %w", err)
	}
	var customConfig = config.NewConfig()
	err = customCfg.Unmarshal(customConfig)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal custom configuration: %w", err)
	}
	return customConfig, nil

}

func LoadService(conf *config.SkyflowConfig) (workflow.WorkflowService, error) {

	var err error
	var dbClient *rdb.DBClient
	if conf.Persistence == nil {
		err = fmt.Errorf("persistence config is missing")
		return nil, err
	}
	dbClient, err = rdb.NewDBClient(conf.Persistence)
	if err != nil {
		return nil, fmt.Errorf("failed to create db client: %w", err)
	}
	if conf.MessageQueue == nil {
		err = fmt.Errorf("message_queue config is missing")
		return nil, err
	}

	innerMQ, err := queue.NewInnerMessageQueueFromConfig(conf.MessageQueue.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create inner message queue: %w", err)
	}

	var delayMQ queue.InnerMessageQueue
	// delay message queue is optional
	if conf.DelayMesageQueue != nil && conf.DelayMesageQueue.DSN != "" {
		delayMQ, err = queue.NewInnerMessageQueueFromConfig(conf.DelayMesageQueue.DSN)
		if err != nil {
			return nil, fmt.Errorf("failed to create delay message queue: %w", err)
		}
	} else {
		delayMQ = queue.NewDBMessageQueue(dbClient, queue.NewDBQueueOption())
	}

	innerMQGroup, err := queue.NewInnerQueueGroup(innerMQ, delayMQ)
	if err != nil {
		return nil, fmt.Errorf("failed to create inner message queue group: %w", err)
	}

	workflowService, err := workflow.NewWorkflowService(dbClient, innerMQGroup, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow service: %w", err)
	}
	return workflowService, nil
}
