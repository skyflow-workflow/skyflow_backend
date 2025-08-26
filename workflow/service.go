package workflow

import (
	"context"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/executor"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/exporter"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/template"
)

// WorkflowService Service provides workflow related services
// It is used to manage workflow execution, templates, and events.
// It is also used to manage workflow state and activity tasks.
// The WorkflowService interface defines the methods that the service should implement.

type WorkflowService = *workflowService

// workflowService is the implementation of WorkflowService interface
// It provides methods to manage workflow execution, templates, and events.
type workflowService struct {
	DBClient         *rdb.DBClient
	TemplateService  template.TemplateService
	InnerQueue       queue.InnerMessageQueue
	Exporter         exporter.ExporterService
	standardExecutor *executor.Executor
	expressExecutor  *executor.Executor
}

func NewWorkflowService(dbClient *rdb.DBClient, innerQueue queue.InnerMessageQueue) (WorkflowService, error) {
	templateService := template.NewTemplateService(dbClient)
	exporterService, err := exporter.NewExporterService(exporter.NewDBListener(dbClient))
	if err != nil {
		return nil, err
	}

	standardExecutor := executor.StandardExecutor
	standardExecutor.MetaDB = dbClient
	standardExecutor.InnerQueue = innerQueue
	standardExecutor.Exporter = exporterService

	expressExecutor := executor.ExpressExecutor
	expressExecutor.MetaDB = dbClient
	expressExecutor.InnerQueue = innerQueue
	expressExecutor.Exporter = exporterService

	svc := &workflowService{
		DBClient:         dbClient,
		TemplateService:  templateService,
		InnerQueue:       innerQueue,
		Exporter:         exporterService,
		standardExecutor: standardExecutor,
		expressExecutor:  expressExecutor,
	}
	return svc, nil
}

func (svc *workflowService) SyncSchema() error {
	var err error
	ctx := context.Background()
	err = svc.TemplateService.SyncSchema(ctx, nil)
	if err != nil {
		return err
	}
	err = svc.Exporter.SyncSchema()
	if err != nil {
		return err
	}
	err = svc.InnerQueue.SyncSchema()
	if err != nil {
		return err
	}
	// TODO: sync executor schema
	return nil
}
