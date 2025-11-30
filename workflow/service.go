package workflow

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backend/config"
	"github.com/skyflow-workflow/skyflow_backend/workflow/executor"
	"github.com/skyflow-workflow/skyflow_backend/workflow/exporter"
	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"
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
	ExporterService  exporter.ExporterService
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
		ExporterService:  exporterService,
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
	err = svc.ExporterService.SyncSchema()
	if err != nil {
		return err
	}
	err = svc.ExecutionService.SyncSchema(ctx, nil)
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

// 创建一个执行
func (svc *workflowService) StopExecution(ctx context.Context, req vo.StopExecutionRequest) (err error) {

	err = svc.ExecutionService.StopExecution(ctx, req)
	if err != nil {
		return
	}
	return
}

// DescribeExecution 获得一个执行的描述
func (svc *workflowService) DescribeExecution(req vo.DescribeExecutionRequest) (dbExe *po.Execution, err error) {

	dbExe, err = svc.ExecutionService.QueryExecutionByUUID(req.ExecutionUUID, executor.ExecutionFields.L5, nil)
	if err != nil {
		return
	}
	return dbExe, nil
}

// DescribeExecutionBone 获得一个执行的框架的描述
func (svc *workflowService) DescribeExecutionBone(ctx context.Context, req vo.DescribeExecutionBoneRequest) (resp vo.DescribeExecutionBoneResponse, err error) {

	resp, err = svc.ExecutionService.DescribeExecutionBone(ctx, req)
	if err != nil {
		return
	}
	return resp, nil
}

// DescribeExecution 获得一个执行的描述
func (svc *workflowService) ListExecutions(req vo.ListExecutionsRequest) (vo.ListExecutionsResponse, error) {

	resp, err := svc.ExecutionService.ListExecutions(req)
	return resp, err
}

// DescribeStep 获得一个执行的描述
func (svc *workflowService) DescribeStep(req vo.DescribeStepRequest) (resp vo.DescribeStepResponse, err error) {

	var dbStep *po.Step
	var dbExe *po.Execution
	dbStep, err = svc.ExecutionService.QueryStepByID(int(req.StepID), executor.StepFields.L5, nil)
	if err != nil {
		return
	}
	dbExe, err = svc.ExecutionService.QueryExecutionByID(dbStep.ExecutionID, executor.ExecutionFields.L1, nil)
	if err != nil {
		return
	}
	resp = vo.DescribeStepResponse{
		Step:          *dbStep,
		ExecutionUUID: dbExe.UUID,
	}

	return resp, nil
}

func (svc *workflowService) ParseWorkflow(req vo.ParseFlowRequest) (*states.StateMachine, error) {
	// Request Limit Check
	smObj, err := svc.standardExecutor.Parser.ParseStateMachine(req.StateMachineDefinition)
	if err != nil {
		return nil, err
	}
	return smObj, err
}

func (svc *workflowService) GetActivityTask(ctx context.Context, req vo.GetActivityTaskRequest) (vo.GetActivityTaskResponse, error) {

	resp, err := svc.ExecutionService.GetActivityTask(ctx, req)
	return resp, err
}

func (svc *workflowService) SendTaskFailure(ctx context.Context, req vo.SendTaskFailureRequest) error {
	// Request Limit Check
	err := svc.ExecutionService.SendTaskFailure(ctx, req)

	return err
}

func (svc *workflowService) SendTaskSuccess(ctx context.Context, req vo.SendTaskSuccessRequest) error {
	// Request Limit Check

	if len(req.Output) > svc.standardExecutor.Config.Quota.MaxTaskOutputSize {
		return vo.ErrorOutputSizeLimitExceeded
	}

	err := svc.ExecutionService.SendTaskSuccess(ctx, req)
	return err
}

func (svc *workflowService) SendTaskReference(ctx context.Context, req vo.SendTaskReferenceRequest) error {
	err := svc.ExecutionService.SendTaskReference(ctx, req)
	return err
}

func (svc *workflowService) SendTaskHeartbeat(ctx context.Context, req vo.SendTaskHeartbeatRequest) error {
	err := svc.ExecutionService.SendTaskHeartbeat(ctx, req)
	return err
}
func (svc *workflowService) SendStepSkip(ctx context.Context, req vo.SendStepSkipRequest) error {
	err := svc.ExecutionService.SendStepSkip(ctx, req)
	return err
}
func (svc *workflowService) RedoStep(ctx context.Context, req vo.DescribeStepRequest) error {
	err := svc.ExecutionService.RedoStep(int(req.StepID))
	return err
}

func (svc *workflowService) ResumeExecution(req vo.DescribeExecutionRequest) error {

	dbExe, err := svc.ExecutionService.QueryExecutionByUUID(req.ExecutionUUID, executor.ExecutionFields.L1, nil)
	if err != nil {
		return err
	}

	err = svc.ExecutionService.ResumeExecution(dbExe.ID)
	if err != nil {
		return err
	}
	return nil
}

func (svc *workflowService) ResumeSuspendingStep(req vo.DescribeStepRequest) error {
	err := svc.ExecutionService.ResumeSuspendingStep(int(req.StepID))
	return err
}

func (svc *workflowService) ListExecutionEvents(req vo.ListExecutionEventsRequest) (vo.ListExecutionEventsResponse, error) {

	dbExe, err := svc.ExecutionService.QueryExecutionByUUID(req.ExecutionUUID, executor.ExecutionFields.L1, nil)
	if err != nil {
		return vo.ListExecutionEventsResponse{}, err
	}
	req.ExecutionID = dbExe.ID
	resp, err := svc.ExporterService.ListExecutionEvents(req)
	return resp, err
}

func (svc *workflowService) ListStepEvents(req vo.ListStepEventsRequest) (vo.ListExecutionEventsResponse, error) {

	resp, err := svc.ExporterService.ListStepEvents(req)
	return resp, err
}

func (svc *workflowService) SendStepRetry(req vo.DescribeStepRequest) error {

	err := svc.ExecutionService.SendStepRetry(int(req.StepID))

	return err
}

func (svc *workflowService) SendStepFailed(ctx context.Context, req vo.DescribeStepRequest) error {

	err := svc.ExecutionService.SendStepFailed(int(req.StepID))

	return err
}

func (svc *workflowService) RetryExecution(ctx context.Context, req vo.DescribeExecutionRequest) error {

	dbExe, err := svc.ExecutionService.QueryExecutionByUUID(req.ExecutionUUID, executor.ExecutionFields.L1, nil)
	if err != nil {
		return err
	}
	err = svc.ExecutionService.RetryExecution(dbExe.ID)
	return err
}

func (svc *workflowService) SkipBlockedTask(ctx context.Context, req vo.SkipBlockedTaskRequest) error {

	err := svc.ExecutionService.SkipBlockedTask(ctx, req)
	return err
}

func (svc *workflowService) UnblockTask(ctx context.Context, req vo.UnblockTaskRequest) error {
	err := svc.ExecutionService.UnblockTask(ctx, req)
	return err
}

func (svc *workflowService) UnblockExecution(ctx context.Context, req vo.DescribeExecutionRequest) error {

	dbExe, err := svc.ExecutionService.QueryExecutionByUUID(req.ExecutionUUID, executor.ExecutionFields.L1, nil)
	if err != nil {
		return err
	}

	err = svc.ExecutionService.UnblockExecution(ctx, dbExe.ID)
	if err != nil {
		return err
	}
	return nil
}
