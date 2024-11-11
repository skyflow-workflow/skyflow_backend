package apiserver

import (
	"github.com/skyflow-workflow/skyflow_backbend/gen/pb"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/server"
)

// ApiServer is the API server.
type ApiServer struct {
	server *server.Server
}

// NewApiServer creates a new API server.
func NewApiServer() *ApiServer {
	s := trpc.NewServer()
	pb.RegisterCommonServiceService(s, &CommonService{})
	pb.RegisterSkyflowServiceService(s, &SkyflowService{})

	return &ApiServer{
		server: s,
	}
}

// Start starts the API server.
func (s *ApiServer) Start() {
	// Start the API server
	if err := s.server.Serve(); err != nil {
		log.Error(err)
	}
}
