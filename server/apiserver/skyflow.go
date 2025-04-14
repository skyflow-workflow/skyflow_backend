package apiserver

import (
	"context"

	"github.com/skyflow-workflow/skyflow_backbend/gen/pb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/template"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

// SkyflowService ...
type SkyflowService struct {
	templateService template.TemplateService
}

// CreateOrUpdateActivity implements pb.SkyflowServiceService.
func (s *SkyflowService) CreateOrUpdateActivity(ctx context.Context, req *pb.CreateActivityRequest) (*pb.CreateActivityResponse, error) {
	panic("unimplemented")
}

// CreateOrUpdateNamespace implements pb.SkyflowServiceService.
func (s *SkyflowService) CreateOrUpdateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest) (*pb.CreateNamespaceResponse, error) {
	panic("unimplemented")
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
		ActivityUri: voresp.URI,
		CreateTime:  voresp.CreateTime.Unix(),
		UpdateTime:  voresp.ModifyTime.Unix(),
	}
	return &resp, nil
}

// CreateNamespace implements pb.SkyflowService.
func (s *SkyflowService) CreateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest) (*pb.CreateNamespaceResponse, error) {
	voreq := vo.CreateActivityResponse
	panic("unimplemented")
}

// DescribeActivity implements pb.SkyflowService.
func (s *SkyflowService) DescribeActivity(ctx context.Context, req *pb.DescribeActivityRequest) (*pb.DescribeActivityResponse, error) {
	panic("unimplemented")
}

// ListActivities implements pb.SkyflowService.
func (s *SkyflowService) ListActivities(ctx context.Context, req *pb.ListActivitiesRequest) (*pb.ListActivitiesResponse, error) {
	panic("unimplemented")
}

// ListNamespaces implements pb.SkyflowService.
func (s *SkyflowService) ListNamespaces(ctx context.Context, req *pb.ListNamespacesRequest) (*pb.ListNamespacesResponse, error) {
	panic("unimplemented")
}
