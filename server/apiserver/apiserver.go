package apiserver

import (
	pbv1 "github.com/skyflow-workflow/skyflow_backbend/api/v1"
	"github.com/skyflow-workflow/skyflow_backbend/workflow"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/server"
)

// APIServer is the API server.
type APIServer struct {
	server      *server.Server
	workflowSvc workflow.WorkflowService
}

// NewAPIServer creates a new API server.
func NewAPIServer(server *server.Server, svc workflow.WorkflowService) *APIServer {
	pbv1.RegisterCommonServiceService(server, &CommonServiceHandler{})
	pbv1.RegisterSkyflowV1ServiceService(server, &SkyflowServiceHandler{
		wfSvc: svc,
	})

	return &APIServer{
		server:      server,
		workflowSvc: svc,
	}
}

// Start starts the API server.
func (s *APIServer) Start() {

	// Start the API server
	if err := s.server.Serve(); err != nil {
		log.Error(err)
	}
}
