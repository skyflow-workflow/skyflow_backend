package apiserver

import (
	"context"
	"log/slog"

	"github.com/skyflow-workflow/skyflow_backbend/gen/pb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

// SkyflowServiceHandler skyflow service handler
type SkyflowServiceHandler struct {
	wfSvc workflow.WorkflowService
}

// CreateOrUpdateStateMachine implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) CreateOrUpdateStateMachine(ctx context.Context, req *pb.CreateStateMachineRequest) (*pb.CreateStateMachineResponse, error) {
	panic("unimplemented")
}

// CreateStateMachine implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) CreateStateMachine(ctx context.Context, req *pb.CreateStateMachineRequest) (*pb.CreateStateMachineResponse, error) {
	panic("unimplemented")
}

// DeleteActivity implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DeleteActivity(ctx context.Context, req *pb.DeleteActivityRequest) (*pb.DeleteActivityResponse, error) {
	panic("unimplemented")
}

// DeleteStateMachine implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DeleteStateMachine(ctx context.Context, req *pb.DeleteStateMachineRequest) (*pb.DeleteStateMachineResponse, error) {
	panic("unimplemented")
}

// DescribeStateMachine implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) DescribeStateMachine(ctx context.Context, req *pb.DescribeStateMachineRequest) (*pb.DescribeStateMachineResponse, error) {
	panic("unimplemented")
}

// ListStateMachines implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) ListStateMachines(ctx context.Context, req *pb.ListStateMachinesRequest) (*pb.ListStateMachinesResponse, error) {
	panic("unimplemented")
}

// UpdateStateMachine implements pb.SkyflowV1ServiceService.
func (s *SkyflowServiceHandler) UpdateStateMachine(ctx context.Context, req *pb.UpdateStateMachineRequest) (*pb.UpdateStateMachineResponse, error) {
	panic("unimplemented")
}

// CreateOrUpdateActivity implements pb.SkyflowServiceService.
func (s *SkyflowServiceHandler) CreateOrUpdateActivity(ctx context.Context, req *pb.CreateActivityRequest) (*pb.CreateActivityResponse, error) {
	voreq := vo.CreateActivityRequest{
		ActivityName: req.Name,
		Comment:      req.Comment,
		Namespace:    req.Namespace,
	}
	voresp, err := s.wfSvc.TemplateService.CreateOrUpdateActivity(ctx, voreq, nil)
	if err != nil {
		return nil, err
	}
	resp := pb.CreateActivityResponse{
		ActivityUri: voresp.Data.URI,
		CreateTime:  voresp.Data.CreateTime.Unix(),
		UpdateTime:  voresp.Data.UpdateTime.Unix(),
	}
	return &resp, nil
}

// CreateOrUpdateNamespace implements pb.SkyflowServiceService.
func (s *SkyflowServiceHandler) CreateOrUpdateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest) (*pb.CreateNamespaceResponse, error) {
	voreq := vo.CreateNamespaceRequest{
		Name:    req.Name,
		Comment: req.Comment,
	}
	voresp, err := s.wfSvc.TemplateService.CreateOrUpdateNamespace(ctx, voreq, nil)
	if err != nil {
		return nil, err
	}
	resp := pb.CreateNamespaceResponse{
		Name:       voresp.Data.Name,
		CreateTime: voresp.Data.CreateTime.Unix(),
		UpdateTime: voresp.Data.UpdateTime.Unix(),
	}
	return &resp, nil
}

// CreateActivity implements pb.SkyflowService.
func (s *SkyflowServiceHandler) CreateActivity(ctx context.Context, req *pb.CreateActivityRequest) (*pb.CreateActivityResponse, error) {
	voreq := vo.CreateActivityRequest{
		ActivityName: req.Name,
		Comment:      req.Comment,
		Namespace:    req.Namespace,
	}
	voresp, err := s.wfSvc.TemplateService.CreateActivity(ctx, voreq, nil)
	if err != nil {
		return nil, err
	}
	resp := pb.CreateActivityResponse{
		ActivityUri: voresp.Data.URI,
		CreateTime:  voresp.Data.CreateTime.Unix(),
		UpdateTime:  voresp.Data.UpdateTime.Unix(),
	}
	return &resp, nil
}

// CreateNamespace implements pb.SkyflowService.
func (s *SkyflowServiceHandler) CreateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest) (*pb.CreateNamespaceResponse, error) {

	voreq := vo.CreateNamespaceRequest{
		Name:    req.Name,
		Comment: req.Comment,
	}
	voresp, err := s.wfSvc.TemplateService.CreateNamespace(ctx, voreq, nil)
	if err != nil {
		return nil, err
	}
	resp := pb.CreateNamespaceResponse{
		Name:       voresp.Data.Name,
		CreateTime: voresp.Data.CreateTime.Unix(),
		UpdateTime: voresp.Data.UpdateTime.Unix(),
	}
	return &resp, nil
}

// DescribeActivity implements pb.SkyflowService.
func (s *SkyflowServiceHandler) DescribeActivity(ctx context.Context, req *pb.DescribeActivityRequest) (*pb.DescribeActivityResponse, error) {
	voresp, err := s.wfSvc.TemplateService.DescribeActivity(ctx, req.ActivityUri, nil)
	if err != nil {
		return nil, err
	}
	resp := pb.DescribeActivityResponse{
		Name:        voresp.Name,
		ActivityUri: voresp.URI,
		Comment:     voresp.Comment,
		CreateTime:  voresp.CreateTime.Unix(),
		UpdateTime:  voresp.UpdateTime.Unix(),
	}
	return &resp, nil
}

// ListActivities implements pb.SkyflowService.
func (s *SkyflowServiceHandler) ListActivities(ctx context.Context, req *pb.ListActivitiesRequest) (*pb.ListActivitiesResponse, error) {

	voreq := vo.ListActivitiesRequest{
		PageRequest: ToVOPageRequest(req.PageRequest),
	}

	voresp, err := s.wfSvc.TemplateService.ListActivities(ctx, voreq)
	if err != nil {
		return nil, err
	}
	respdata := DataTransferArray(voresp.Activities, ToPBActivity)

	resp := &pb.ListActivitiesResponse{
		Activities:   respdata,
		PageResponse: ToPBPageResponse(voresp.PageResponse),
	}
	return resp, nil

}

// ListNamespaces implements pb.SkyflowService.
func (s *SkyflowServiceHandler) ListNamespaces(ctx context.Context, req *pb.ListNamespacesRequest) (*pb.ListNamespacesResponse, error) {

	slog.Info("ListNamespaces called", "req", req)
	voreq := vo.ListNamespacesRequest{
		PageRequest: ToVOPageRequest(req.PageRequest),
	}

	voresp, err := s.wfSvc.TemplateService.ListNamespaces(ctx, voreq)
	if err != nil {
		return nil, err
	}
	respdata := DataTransferArray(voresp.Namespaces, ToPBNamespace)

	resp := &pb.ListNamespacesResponse{
		Namespaces:   respdata,
		PageResponse: ToPBPageResponse(voresp.PageResponse),
	}
	return resp, nil
}
