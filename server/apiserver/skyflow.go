package apiserver

import (
	"context"

	"github.com/skyflow-workflow/skyflow_backbend/gen/pb"
)

type SkyflowService struct{}

// CreateActivity implements pb.SkyflowService.
func (s *SkyflowService) CreateActivity(ctx context.Context, req *pb.CreateActivityRequest) (*pb.CreateActivityResponse, error) {
	panic("unimplemented")
}

// CreateNamespace implements pb.SkyflowService.
func (s *SkyflowService) CreateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest) (*pb.CreateNamespaceResponse, error) {
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
