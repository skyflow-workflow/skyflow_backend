package workflow

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backend/config"
	"github.com/skyflow-workflow/skyflow_backend/workflow/executor"
	"github.com/skyflow-workflow/skyflow_backend/workflow/exporter"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backend/workflow/template"
	"github.com/skyflow-workflow/skyflow_backend/workflow/vo"
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
	ExecutionService executor.ExecutionService
	standardExecutor *executor.Executor
	expressExecutor  *executor.Executor
}

func NewWorkflowService(dbClient *rdb.DBClient, innerQueue queue.InnerMessageQueue, conf *config.SkyflowConfig) (WorkflowService, error) {
	templateService := template.NewTemplateService(dbClient)
	exporterService, err := exporter.NewExporterService(dbClient, conf.Exporter)
	if err != nil {
		return nil, err
	}
	executionService := executor.NewExecutionService(dbClient, innerQueue, exporterService)

	standardExecutor := executor.StandardExecutor

	standardExecutor.ExecutionService = executionService

	expressExecutor := executor.ExpressExecutor
	expressExecutor.ExecutionService = executionService

	svc := &workflowService{
		DBClient:         dbClient,
		TemplateService:  templateService,
		InnerQueue:       innerQueue,
		Exporter:         exporterService,
		ExecutionService: executionService,
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
	err = svc.ExecutionService.SyncSchema(ctx, nil)
	if err != nil {
		return err
	}
	err = svc.InnerQueue.SyncSchema()
	if err != nil {
		return err
	}

	return nil
}

func (svc *workflowService) ValidateStartExecutionRequest(ctx context.Context, req *vo.StartExecutionRequest) (err error) {

	var dbSM po.StateMachine
	// Request Limit Check
	if len(req.Input) > svc.ExecutionService.StandardExecutor.Config.Quota.MaxInputSize {
		err = vo.ErrorStartExecutionInputSizeLimitExceeded
		return
	}
	if len(req.Title) > svc.ExecutionService.StandardExecutor.Config.Quota.MaxTitleSize {
		err = vo.ErrorTitleSizeLimitExceeded
		return
	}
	if len(req.ExecutionUUID) > svc.ExecutionService.StandardExecutor.Config.Quota.MaxExecutionUUIDSize {
		err = vo.ErrorExecutionUUIDSizeLimitExceed
		return
	}

	// 优先看definition
	if req.StateMachineDefinition != "" {
		// Request Limit Check
		if len(req.StateMachineDefinition) > svc.ExecutionService.StandardExecutor.Config.Quota.MaxStateMachineSize {
			err = vo.ErrorStateMachineSizeLimitExceeded
			return
		}
		req.StateMachineURI = ""
	} else if req.StateMachineURI != "" {
		// 再看URI
		// URI 长度验证通过proto validate
		dbSM, err = svc.TemplateService.DescribeStateMachine(ctx,
			vo.DescribeStateMachineRequest{
				StateMachineURI: req.StateMachineURI,
			}, nil)
		if err != nil {
			return
		}
		req.StateMachineDefinition = dbSM.Definition
	} else {
		err = fmt.Errorf("%w: must indicate statemachine_defintion/statemachine_uri", vo.ErrorParameterInvalid)
		return
	}
	return
}

func (svc *workflowService) StartExecution(ctx context.Context, req vo.StartExecutionRequest) (dbExe *po.Execution, err error) {
	slog.Info("start execution ", req)
	err = svc.ValidateStartExecutionRequest(ctx, &req)
	if err != nil {
		return
	}
	dbExe, err = svc.ExecutionService.StartExecution(req)
	if err != nil {
		return
	}
	return
}
