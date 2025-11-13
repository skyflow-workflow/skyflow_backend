package apiserver

import (
	"context"

	"github.com/skyflow-workflow/skyflow_backbend/gen/pb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SkyflowServiceHandler skyflow service handler
type SkyflowServiceHandler struct {
	wfSvc workflow.WorkflowService
}

// ValidateStateMachineDefinition implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) ValidateStateMachineDefinition(ctx context.Context, req *pb.ValidateStateMachineDefinitionRequest) (*pb.ValidateStateMachineDefinitionResponse, error) {
	panic("unimplemented")
}

// DescribeExecution implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DescribeExecution(ctx context.Context, req *pb.DescribeExecutionRequest) (*pb.DescribeExecutionResponse, error) {
	panic("unimplemented")
}

// DescribeExecutionBone implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DescribeExecutionBone(ctx context.Context, req *pb.DescribeExecutionBoneRequest) (*pb.DescribeExecutionBoneResponse, error) {
	panic("unimplemented")
}

// DescribeStep implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DescribeStep(ctx context.Context, req *pb.DescribeStepRequest) (*pb.DescribeStepResponse, error) {
	panic("unimplemented")
}

// GetActivityTask implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) GetActivityTask(ctx context.Context, req *pb.GetActivityTaskRequest) (*pb.GetActivityTaskResponse, error) {
	panic("unimplemented")
}

// ListExecutionEvents implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) ListExecutionEvents(ctx context.Context, req *pb.ListExecutionEventsRequest) (*pb.ListExecutionEventsResponse, error) {
	panic("unimplemented")
}

// ListExecutions implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) ListExecutions(ctx context.Context, req *pb.ListExecutionsRequest) (*pb.ListExecutionsResponse, error) {
	panic("unimplemented")
}

// ListStepEvents implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) ListStepEvents(ctx context.Context, req *pb.ListStepEventsRequest) (*pb.ListExecutionEventsResponse, error) {
	panic("unimplemented")
}

// ParseStateMachine implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) ParseStateMachine(ctx context.Context, req *pb.ParseStateMachineRequest) (*pb.ParseStateMachineResponse, error) {
	panic("unimplemented")
}

// StartExecution implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) StartExecution(ctx context.Context, req *pb.StartExecutionRequest) (*pb.StartExecutionResponse, error) {
	voReq := vo.StartExecutionRequest{
		StateMachineURI:        req.StatemachineUri,
		Input:                  req.Input,
		StateMachineDefinition: req.Definition,
		Title:                  req.Title,
		ExecutionUUID:          req.ExecutionName,
	}
	dbExecution, err := s.wfSvc.ExecutionService.StartExecution(voReq)
	if err != nil {
		return nil, err
	}
	resp := &pb.StartExecutionResponse{
		ExecutionUuid: dbExecution.UUID,
		CreateTime:    dbExecution.CreateTime.String(),
	}
	return resp, nil
}

// StopExecution implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) StopExecution(ctx context.Context, req *pb.StopExecutionRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// DeleteNamespace implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DeleteNamespace(ctx context.Context, req *pb.DeleteNamespaceRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// CreateOrUpdateStateMachine implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) CreateOrUpdateStateMachine(ctx context.Context, req *pb.CreateStateMachineRequest) (*pb.CreateStateMachineResponse, error) {
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
	resp := &pb.CreateStateMachineResponse{
		Data: ToPBStateMachine(&voResp.Data),
	}
	return resp, nil
}

// CreateStateMachine implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) CreateStateMachine(ctx context.Context, req *pb.CreateStateMachineRequest) (*pb.CreateStateMachineResponse, error) {
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
	resp := &pb.CreateStateMachineResponse{
		Data: ToPBStateMachine(&voResp.Data),
	}
	return resp, nil
}

// DeleteActivity implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DeleteActivity(ctx context.Context, req *pb.DeleteActivityRequest) (*pb.DeleteActivityResponse, error) {
	voReq := vo.DeleteActivityRequest{
		ActivityURI: req.ActivityUri,
	}
	err := s.wfSvc.TemplateService.DeleteActivity(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := &pb.DeleteActivityResponse{}
	return resp, nil
}

// DeleteStateMachine implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DeleteStateMachine(ctx context.Context, req *pb.DeleteStateMachineRequest) (*pb.DeleteStateMachineResponse, error) {
	voReq := vo.DeleteStateMachineRequest{
		StateMachineURI: req.StatemachineUri,
	}
	err := s.wfSvc.TemplateService.DeleteStateMachine(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteStateMachineResponse{}, nil

}

// DescribeStateMachine implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DescribeStateMachine(ctx context.Context, req *pb.DescribeStateMachineRequest) (*pb.DescribeStateMachineResponse, error) {
	voReq := vo.DescribeStateMachineRequest{
		StateMachineURI: req.StatemachineUri,
	}
	voResp, err := s.wfSvc.TemplateService.DescribeStateMachine(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	return &pb.DescribeStateMachineResponse{
		Data: ToPBStateMachine(&voResp),
	}, nil
}

// ListStateMachines implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) ListStateMachines(ctx context.Context, req *pb.ListStateMachinesRequest) (*pb.ListStateMachinesResponse, error) {

	voReq := vo.ListStateMachinesRequest{
		Namespace:   req.Namespace,
		PageRequest: ToVOPageRequest(req.PageRequest),
	}
	voResp, err := s.wfSvc.TemplateService.ListStateMachines(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := &pb.ListStateMachinesResponse{
		Statemachines: DataTransferArray(voResp.StateMachines, ToPBStateMachineItem),
		PageResponse:  ToPBPageResponse(voResp.PageResponse),
	}
	return resp, nil
}

// UpdateStateMachine implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) UpdateStateMachine(ctx context.Context, req *pb.UpdateStateMachineRequest) (*pb.UpdateStateMachineResponse, error) {
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
	resp := &pb.UpdateStateMachineResponse{}
	return resp, nil
}

// CreateOrUpdateActivity implements pb.SkyflowServiceService.
func (s *SkyflowServiceHandler) CreateOrUpdateActivity(ctx context.Context, req *pb.CreateActivityRequest) (*pb.CreateActivityResponse, error) {
	voReq := vo.CreateActivityRequest{
		ActivityName: req.Name,
		Description:  req.Description,
		Namespace:    req.Namespace,
	}
	voResp, err := s.wfSvc.TemplateService.CreateOrUpdateActivity(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := pb.CreateActivityResponse{
		ActivityUri: voResp.Data.URI,
		CreateTime:  voResp.Data.CreateTime.Unix(),
		UpdateTime:  voResp.Data.UpdateTime.Unix(),
	}
	return &resp, nil
}

// CreateOrUpdateNamespace implements pb.SkyflowServiceService.
func (s *SkyflowServiceHandler) CreateOrUpdateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest) (*pb.CreateNamespaceResponse, error) {
	voReq := vo.CreateNamespaceRequest{
		Name:        req.Name,
		Description: req.Description,
	}
	voResp, err := s.wfSvc.TemplateService.CreateOrUpdateNamespace(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := pb.CreateNamespaceResponse{
		Name:       voResp.Data.Name,
		CreateTime: voResp.Data.CreateTime.Unix(),
		UpdateTime: voResp.Data.UpdateTime.Unix(),
	}
	return &resp, nil
}

// CreateActivity implements pb.SkyflowService.
func (s *SkyflowServiceHandler) CreateActivity(ctx context.Context, req *pb.CreateActivityRequest) (*pb.CreateActivityResponse, error) {
	voReq := vo.CreateActivityRequest{
		ActivityName: req.Name,
		Description:  req.Description,
		Namespace:    req.Namespace,
	}
	voResp, err := s.wfSvc.TemplateService.CreateActivity(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := pb.CreateActivityResponse{
		ActivityUri: voResp.Data.URI,
		CreateTime:  voResp.Data.CreateTime.Unix(),
		UpdateTime:  voResp.Data.UpdateTime.Unix(),
	}
	return &resp, nil
}

// CreateNamespace implements pb.SkyflowService.
func (s *SkyflowServiceHandler) CreateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest) (*pb.CreateNamespaceResponse, error) {

	voReq := vo.CreateNamespaceRequest{
		Name:        req.Name,
		Description: req.Description,
	}
	voResp, err := s.wfSvc.TemplateService.CreateNamespace(ctx, voReq, nil)
	if err != nil {
		return nil, err
	}
	resp := pb.CreateNamespaceResponse{
		Name:       voResp.Data.Name,
		CreateTime: voResp.Data.CreateTime.Unix(),
		UpdateTime: voResp.Data.UpdateTime.Unix(),
	}
	return &resp, nil
}

// DescribeActivity implements pb.SkyflowService.
func (s *SkyflowServiceHandler) DescribeActivity(ctx context.Context, req *pb.DescribeActivityRequest) (*pb.DescribeActivityResponse, error) {
	voResp, err := s.wfSvc.TemplateService.DescribeActivity(ctx, req.ActivityUri, nil)
	if err != nil {
		return nil, err
	}
	resp := pb.DescribeActivityResponse{
		Name:        voResp.Name,
		ActivityUri: voResp.URI,
		Description: voResp.Description,
		CreateTime:  voResp.CreateTime.Unix(),
		UpdateTime:  voResp.UpdateTime.Unix(),
	}
	return &resp, nil
}

// ListActivities implements pb.SkyflowService.
func (s *SkyflowServiceHandler) ListActivities(ctx context.Context, req *pb.ListActivitiesRequest) (*pb.ListActivitiesResponse, error) {

	voReq := vo.ListActivitiesRequest{
		PageRequest: ToVOPageRequest(req.PageRequest),
	}

	voResp, err := s.wfSvc.TemplateService.ListActivities(ctx, voReq)
	if err != nil {
		return nil, err
	}
	respdata := DataTransferArray(voResp.Activities, ToPBActivityItem)

	resp := &pb.ListActivitiesResponse{
		Activities:   respdata,
		PageResponse: ToPBPageResponse(voResp.PageResponse),
	}
	return resp, nil

}

// ListNamespaces implements pb.SkyflowService.
func (s *SkyflowServiceHandler) ListNamespaces(ctx context.Context, req *pb.ListNamespacesRequest) (*pb.ListNamespacesResponse, error) {

	voReq := vo.ListNamespacesRequest{
		PageRequest: ToVOPageRequest(req.PageRequest),
	}

	voResp, err := s.wfSvc.TemplateService.ListNamespaces(ctx, voReq)
	if err != nil {
		return nil, err
	}
	respdata := DataTransferArray(voResp.Namespaces, ToPBNamespace)

	resp := &pb.ListNamespacesResponse{
		Namespaces:   respdata,
		PageResponse: ToPBPageResponse(voResp.PageResponse),
	}
	return resp, nil
}
