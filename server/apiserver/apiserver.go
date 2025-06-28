package apiserver

import (
	pb "github.com/skyflow-workflow/skyflow_backbend/gen/pb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/server"
)

// APIServer is the API server.
type APIServer struct {
	server      *server.Server
	workflowSvc workflow.WorkflowService
}

var skyflowConfigFilePath string = "./skyflow.yaml"

// NewAPIServer creates a new API server.
func NewAPIServer(svc workflow.WorkflowService) *APIServer {

	s := trpc.NewServer()
	pb.RegisterCommonServiceService(s, &CommonServiceHandler{})
	pb.RegisterSkyflowV1ServiceService(s, &SkyflowServiceHandler{
		wfSvc: svc,
	})

	return &APIServer{
		server:      s,
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

func LoadConfig(configFilePath string) error {
	// Load the TRPC server configuration
	if configFilePath != "" {
		trpc.ServerConfigPath = configFilePath
	}
	_, err := trpc.LoadConfig(configFilePath)
	if err != nil {
		log.Error("Error loading TRPC server configuration", "error", err, "configPath", configFilePath)
		return err
	}
	return nil
}
