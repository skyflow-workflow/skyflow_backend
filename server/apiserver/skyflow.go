package apiserver

import (
	"context"

	"github.com/skyflow-workflow/skyflow_backbend/gen/pb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/template"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

// SkyflowService skyflow service handler
type SkyflowService struct {
	templateService template.TemplateService
}

// CreateOrUpdateActivity implements pb.SkyflowServiceService.
func (s *SkyflowService) CreateOrUpdateActivity(ctx context.Context, req *pb.CreateActivityRequest) (*pb.CreateActivityResponse, error) {
	voreq := vo.CreateActivityRequest{
		ActivityName: req.Name,
		Comment:      req.Comment,
		Namespace:    req.Namespace,
	}
	voresp, err := s.templateService.CreateOrUpdateActivity(ctx, voreq, nil)
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
func (s *SkyflowService) CreateOrUpdateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest) (*pb.CreateNamespaceResponse, error) {
	voreq := vo.CreateNamespaceRequest{
		Name:    req.Name,
		Comment: req.Comment,
	}
	voresp, err := s.templateService.CreateOrUpdateNamespace(ctx, voreq, nil)
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
func (s *SkyflowService) CreateActivity(ctx context.Context, req *pb.CreateActivityRequest) (*pb.CreateActivityResponse, error) {
	voreq := vo.CreateActivityRequest{
		ActivityName: req.Name,
		Comment:      req.Comment,
		Namespace:    req.Namespace,
	}
	voresp, err := s.templateService.CreateActivity(ctx, voreq, nil)
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
func (s *SkyflowService) CreateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest) (*pb.CreateNamespaceResponse, error) {

	voreq := vo.CreateNamespaceRequest{
		Name:    req.Name,
		Comment: req.Comment,
	}
	voresp, err := s.templateService.CreateNamespace(ctx, voreq, nil)
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
func (s *SkyflowService) DescribeActivity(ctx context.Context, req *pb.DescribeActivityRequest) (*pb.DescribeActivityResponse, error) {
	voresp, err := s.templateService.DescribeActivity(ctx, req.ActivityUri, nil)
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
func (s *SkyflowService) ListActivities(ctx context.Context, req *pb.ListActivitiesRequest) (*pb.ListActivitiesResponse, error) {

	voreq := vo.ListActivitiesRequest{
		PageRequest: ToVOPageRequest(req.PageRequest),
	}

	voresp, err := s.templateService.ListActivities(ctx, voreq)
	if err != nil {
		return nil, err
	}
	respdata := DataTransferArray[po.Activity, *pb.ActivityListItem](voresp.Activities, ToPBActivity)

	resp := &pb.ListActivitiesResponse{
		Activities:   respdata,
		PageResponse: ToPBPageResponse(voresp.PageResponse),
	}
	return resp, nil

}

// ListNamespaces implements pb.SkyflowService.
func (s *SkyflowService) ListNamespaces(ctx context.Context, req *pb.ListNamespacesRequest) (*pb.ListNamespacesResponse, error) {

	voreq := vo.ListNamespacesRequest{
		PageRequest: ToVOPageRequest(req.PageRequest),
	}

	voresp, err := s.templateService.ListNamespaces(ctx, voreq)
	if err != nil {
		return nil, err
	}
	respdata := DataTransferArray[po.Namespace, *pb.NamespaceListItem](voresp.Namespaces, ToPBNamespace)

	resp := &pb.ListNamespacesResponse{
		Namespaces:   respdata,
		PageResponse: ToPBPageResponse(voresp.PageResponse),
	}
	return resp, nil
}
