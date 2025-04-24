package apiserver

import (
	"context"

	"github.com/skyflow-workflow/skyflow_backbend/gen/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CommonService is a service that provides common functions.
type CommonService struct {
}

// Ping implements pb.CommonServiceService.
func (c *CommonService) Ping(ctx context.Context, req *emptypb.Empty) (*pb.PingResponse, error) {
	resp := &pb.PingResponse{
		Message: "Pong",
	}
	return resp, nil
}

// HTTP implements pb.CommonService.
func (c *CommonService) HTTP(ctx context.Context, req *emptypb.Empty) (*pb.HTTPResponseMessage, error) {
	resp := &pb.HTTPResponseMessage{
		Retcode: 0,
		Message: "Hello, world!",
	}
	return resp, nil
}

// Paging implements pb.CommonService.
func (c *CommonService) Paging(ctx context.Context, req *pb.PageRequest) (*pb.PageResponse, error) {

	resp := &pb.PageResponse{
		Count:      100,
		PageSize:   10,
		PageNumber: 100,
		PageCount:  10,
	}
	return resp, nil
}
