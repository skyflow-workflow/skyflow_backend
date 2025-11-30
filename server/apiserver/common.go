package apiserver

import (
	"context"

	pbv1 "github.com/skyflow-workflow/skyflow_backend/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CommonServiceHandler is a service that provides common functions.
type CommonServiceHandler struct {
	pbv1.UnimplementedCommonServiceServer
}

// Ping implements pb.CommonServiceService.
func (c *CommonServiceHandler) Ping(ctx context.Context, req *emptypb.Empty) (*pbv1.PingResponse, error) {
	resp := &pbv1.PingResponse{
		Message: "Pong",
	}
	return resp, nil
}

// HTTP implements pb.CommonService.
func (c *CommonServiceHandler) HTTP(ctx context.Context, req *emptypb.Empty) (*pbv1.HTTPResponseMessage, error) {
	resp := &pbv1.HTTPResponseMessage{
		Retcode: 0,
		Message: "Hello, world!",
	}
	return resp, nil
}

// Paging implements pb.CommonService.
func (c *CommonServiceHandler) Paging(ctx context.Context, req *pbv1.PageRequest) (*pbv1.PageResponse, error) {

	resp := &pbv1.PageResponse{
		Count:      100,
		PageSize:   10,
		PageNumber: 100,
		PageCount:  10,
	}
	return resp, nil
}
