package apiserver

import (
	"context"
	"errors"

	pbv1 "github.com/skyflow-workflow/skyflow_backend/api/v1"
	"github.com/skyflow-workflow/skyflow_backend/workflow"
	"github.com/skyflow-workflow/skyflow_backend/workflow/vo"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SkyflowServiceHandler skyflow service handler
type SkyflowServiceHandler struct {
	pbv1.UnimplementedSkyflowV1ServiceServer
	wfSvc workflow.WorkflowService
}

// NewSkyflowServiceHandler creates a new SkyflowServiceHandler
func NewSkyflowServiceHandler(wfSvc workflow.WorkflowService) *SkyflowServiceHandler {
	return &SkyflowServiceHandler{
		wfSvc: wfSvc,
	}
}

// ValidateStateMachineDefinition implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) ValidateStateMachineDefinition(ctx context.Context, req *pbv1.ValidateStateMachineDefinitionRequest) (*pbv1.ValidateStateMachineDefinitionResponse, error) {
	panic("unimplemented")
}

// DescribeExecution implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DescribeExecution(ctx context.Context, req *pbv1.DescribeExecutionRequest) (*pbv1.DescribeExecutionResponse, error) {
	panic("unimplemented")
}

// DescribeExecutionBone implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DescribeExecutionBone(ctx context.Context, req *pbv1.DescribeExecutionBoneRequest) (*pbv1.DescribeExecutionBoneResponse, error) {
	panic("unimplemented")
}

// DescribeStep implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DescribeStep(ctx context.Context, req *pbv1.DescribeStepRequest) (*pbv1.DescribeStepResponse, error) {

	panic("unimplemented")
}

// GetActivityTask implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) GetActivityTask(ctx context.Context, req *pbv1.GetActivityTaskRequest) (*pbv1.GetActivityTaskResponse, error) {

	if req.ActivityUri == "" {
		return nil, errors.New("activity uri is required")
	}
	voReq := DefaultGetActivityRequest
	// 如果没有设置超时时间，则使用默认的超时时间
	voReq.ActivityURI = req.ActivityUri

	voResp, err := s.wfSvc.ExecutionService.GetActivityTask(ctx, voReq)
	if err != nil {
		return nil, err
	}
	resp := &pbv1.GetActivityTaskResponse{
		ActivityUri:      voResp.Resource,
		TaskToken:        voResp.TaskToken,
		Input:            voResp.Input,
		TimeoutSeconds:   int64(voResp.TimeoutSeconds),
		HeartbeatSeconds: int64(voResp.HeartbeatSeconds),
	}
	return resp, nil
}

// ListExecutionEvents implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) ListExecutionEvents(ctx context.Context, req *pbv1.ListExecutionEventsRequest) (*pbv1.ListExecutionEventsResponse, error) {
	panic("unimplemented")
}

// ListExecutions implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) ListExecutions(ctx context.Context, req *pbv1.ListExecutionsRequest) (*pbv1.ListExecutionsResponse, error) {
	panic("unimplemented")
}

// ListStepEvents implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) ListStepEvents(ctx context.Context, req *pbv1.ListStepEventsRequest) (*pbv1.ListExecutionEventsResponse, error) {
	panic("unimplemented")
}

// ParseStateMachine implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) ParseStateMachine(ctx context.Context, req *pbv1.ParseStateMachineRequest) (*pbv1.ParseStateMachineResponse, error) {

	panic("unimplemented")
}

// StartExecution implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) StartExecution(ctx context.Context, req *pbv1.StartExecutionRequest) (*pbv1.StartExecutionResponse, error) {
	// 输入验证
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}
	if req.StatemachineUri == "" && req.Definition == "" {
		return nil, errors.New("state machine definition is required")
	}
	if req.Input == "" {
		req.Input = "{}"
	}

	voReq := vo.StartExecutionRequest{
		StateMachineURI:        req.StatemachineUri,
		Input:                  req.Input,
		StateMachineDefinition: req.Definition,
		Title:                  req.Title,
		ExecutionUUID:          req.ExecutionName,
	}
	dbExecution, err := s.wfSvc.StartExecution(ctx, voReq)
	if err != nil {
		return nil, err
	}
	resp := &pbv1.StartExecutionResponse{
		ExecutionUuid: dbExecution.UUID,
		CreateTime:    dbExecution.CreateTime.String(),
	}
	return resp, nil
}

// StopExecution implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) StopExecution(ctx context.Context, req *pbv1.StopExecutionRequest) (*emptypb.Empty, error) {
	voReq := vo.StopExecutionRequest{
		ExecutionUUID: req.ExecutionUuid,
	}
	err := s.wfSvc.ExecutionService.StopExecution(ctx, voReq)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteNamespace implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DeleteNamespace(ctx context.Context, req *pbv1.DeleteNamespaceRequest) (*emptypb.Empty, error) {
	voReq := vo.DeleteNamespaceRequest{
		Name: req.Name,
	}
	err := s.wfSvc.TemplateService.DeleteNamespace(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// CreateOrUpdateStateMachine implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) CreateOrUpdateStateMachine(ctx context.Context, req *pbv1.CreateStateMachineRequest) (*pbv1.CreateStateMachineResponse, error) {
	voReq := vo.CreateStateMachineRequest{
		Name:        req.Name,
		Description: req.Description,
		Definition:  req.Definition,
		Namespace:   req.Namespace,
	}
	voResp, err := s.wfSvc.TemplateService.CreateOrUpdateStateMachine(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := &pbv1.CreateStateMachineResponse{
		Data: ToPBStateMachine(&voResp.Data),
	}
	return resp, nil
}

// CreateStateMachine implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) CreateStateMachine(ctx context.Context, req *pbv1.CreateStateMachineRequest) (*pbv1.CreateStateMachineResponse, error) {
	voReq := vo.CreateStateMachineRequest{
		Name:        req.Name,
		Description: req.Description,
		Definition:  req.Definition,
		Namespace:   req.Namespace,
	}
	voResp, err := s.wfSvc.TemplateService.CreateStateMachine(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := &pbv1.CreateStateMachineResponse{
		Data: ToPBStateMachine(&voResp.Data),
	}
	return resp, nil
}

// DeleteActivity implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DeleteActivity(ctx context.Context, req *pbv1.DeleteActivityRequest) (*pbv1.DeleteActivityResponse, error) {
	voReq := vo.DeleteActivityRequest{
		ActivityURI: req.ActivityUri,
	}
	err := s.wfSvc.TemplateService.DeleteActivity(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := &pbv1.DeleteActivityResponse{}
	return resp, nil
}

// DeleteStateMachine implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DeleteStateMachine(ctx context.Context, req *pbv1.DeleteStateMachineRequest) (*pbv1.DeleteStateMachineResponse, error) {
	voReq := vo.DeleteStateMachineRequest{
		StateMachineURI: req.StatemachineUri,
	}
	err := s.wfSvc.TemplateService.DeleteStateMachine(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	return &pbv1.DeleteStateMachineResponse{}, nil

}

// DescribeStateMachine implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DescribeStateMachine(ctx context.Context, req *pbv1.DescribeStateMachineRequest) (*pbv1.DescribeStateMachineResponse, error) {
	voReq := vo.DescribeStateMachineRequest{
		StateMachineURI: req.StatemachineUri,
	}
	voResp, err := s.wfSvc.TemplateService.DescribeStateMachine(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	return &pbv1.DescribeStateMachineResponse{
		Data: ToPBStateMachine(&voResp),
	}, nil
}

// ListStateMachines implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) ListStateMachines(ctx context.Context, req *pbv1.ListStateMachinesRequest) (*pbv1.ListStateMachinesResponse, error) {

	voReq := vo.ListStateMachinesRequest{
		Namespace:   req.Namespace,
		PageRequest: ToVOPageRequest(req.PageRequest),
	}
	voResp, err := s.wfSvc.TemplateService.ListStateMachines(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := &pbv1.ListStateMachinesResponse{
		Statemachines: DataTransferArray(voResp.StateMachines, ToPBStateMachineItem),
		PageResponse:  ToPBPageResponse(voResp.PageResponse),
	}
	return resp, nil
}

// UpdateStateMachine implements pbv1.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) UpdateStateMachine(ctx context.Context, req *pbv1.UpdateStateMachineRequest) (*pbv1.UpdateStateMachineResponse, error) {
	voReq := vo.UpdateStateMachineRequest{
		StateMachineURI: req.StatemachineUri,
		Name:            req.Name,
		Description:     req.Description,
		Definition:      req.Definition,
	}
	err := s.wfSvc.TemplateService.UpdateStateMachine(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := &pbv1.UpdateStateMachineResponse{}
	return resp, nil
}

// CreateOrUpdateActivity implements pbv1.SkyflowServiceService.
func (s *SkyflowServiceHandler) CreateOrUpdateActivity(ctx context.Context, req *pbv1.CreateActivityRequest) (*pbv1.CreateActivityResponse, error) {
	voReq := vo.CreateActivityRequest{
		ActivityName: req.Name,
		Description:  req.Description,
		Namespace:    req.Namespace,
	}
	voResp, err := s.wfSvc.TemplateService.CreateOrUpdateActivity(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := pbv1.CreateActivityResponse{
		ActivityUri: voResp.Data.URI,
		CreateTime:  voResp.Data.CreateTime.Unix(),
		UpdateTime:  voResp.Data.UpdateTime.Unix(),
	}
	return &resp, nil
}

// CreateOrUpdateNamespace implements pbv1.SkyflowServiceService.
func (s *SkyflowServiceHandler) CreateOrUpdateNamespace(ctx context.Context, req *pbv1.CreateNamespaceRequest) (*pbv1.CreateNamespaceResponse, error) {
	voReq := vo.CreateNamespaceRequest{
		Name:        req.Name,
		Description: req.Description,
	}
	voResp, err := s.wfSvc.TemplateService.CreateOrUpdateNamespace(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := pbv1.CreateNamespaceResponse{
		Name:        voResp.Data.Name,
		Description: voResp.Data.Description,
		CreateTime:  voResp.Data.CreateTime.Unix(),
		UpdateTime:  voResp.Data.UpdateTime.Unix(),
	}
	return &resp, nil
}

// CreateActivity implements pbv1.SkyflowService.
func (s *SkyflowServiceHandler) CreateActivity(ctx context.Context, req *pbv1.CreateActivityRequest) (*pbv1.CreateActivityResponse, error) {
	voReq := vo.CreateActivityRequest{
		ActivityName: req.Name,
		Description:  req.Description,
		Namespace:    req.Namespace,
	}
	voResp, err := s.wfSvc.TemplateService.CreateActivity(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := pbv1.CreateActivityResponse{
		ActivityUri: voResp.Data.URI,
		CreateTime:  voResp.Data.CreateTime.Unix(),
		UpdateTime:  voResp.Data.UpdateTime.Unix(),
	}
	return &resp, nil
}

// CreateNamespace implements pbv1.SkyflowService.
func (s *SkyflowServiceHandler) CreateNamespace(ctx context.Context, req *pbv1.CreateNamespaceRequest) (*pbv1.CreateNamespaceResponse, error) {

	voReq := vo.CreateNamespaceRequest{
		Name:        req.Name,
		Description: req.Description,
	}
	voResp, err := s.wfSvc.TemplateService.CreateNamespace(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := pbv1.CreateNamespaceResponse{
		Name:        voResp.Data.Name,
		Description: voResp.Data.Description,
		CreateTime:  voResp.Data.CreateTime.Unix(),
		UpdateTime:  voResp.Data.UpdateTime.Unix(),
	}
	return &resp, nil
}

// DescribeActivity implements pbv1.SkyflowService.
func (s *SkyflowServiceHandler) DescribeActivity(ctx context.Context, req *pbv1.DescribeActivityRequest) (*pbv1.DescribeActivityResponse, error) {
	voResp, err := s.wfSvc.TemplateService.DescribeActivity(ctx, req.ActivityUri, nil)
	if err != nil {
		return nil, err
	}
	resp := pbv1.DescribeActivityResponse{
		Name:        voResp.Name,
		ActivityUri: voResp.URI,
		Description: voResp.Description,
		CreateTime:  voResp.CreateTime.Unix(),
		UpdateTime:  voResp.UpdateTime.Unix(),
	}
	return &resp, nil
}

// ListActivities implements pbv1.SkyflowService.
func (s *SkyflowServiceHandler) ListActivities(ctx context.Context, req *pbv1.ListActivitiesRequest) (*pbv1.ListActivitiesResponse, error) {

	voReq := vo.ListActivitiesRequest{
		PageRequest: ToVOPageRequest(req.PageRequest),
	}

	voResp, err := s.wfSvc.TemplateService.ListActivities(ctx, voReq)
	if err != nil {
		return nil, err
	}
	respdata := DataTransferArray(voResp.Activities, ToPBActivityItem)

	resp := &pbv1.ListActivitiesResponse{
		Activities:   respdata,
		PageResponse: ToPBPageResponse(voResp.PageResponse),
	}
	return resp, nil

}

// ListNamespaces implements pbv1.SkyflowService.
func (s *SkyflowServiceHandler) ListNamespaces(ctx context.Context, req *pbv1.ListNamespacesRequest) (*pbv1.ListNamespacesResponse, error) {

	voReq := vo.ListNamespacesRequest{
		PageRequest: ToVOPageRequest(req.PageRequest),
	}

	voResp, err := s.wfSvc.TemplateService.ListNamespaces(ctx, voReq)
	if err != nil {
		return nil, err
	}
	respdata := DataTransferArray(voResp.Namespaces, ToPBNamespace)

	resp := &pbv1.ListNamespacesResponse{
		Namespaces:   respdata,
		PageResponse: ToPBPageResponse(voResp.PageResponse),
	}
	return resp, nil
}
