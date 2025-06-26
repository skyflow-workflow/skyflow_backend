package apiserver

import (
	pb "github.com/skyflow-workflow/skyflow_backbend/gen/pb"
	trpc "trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/server"
)

// ApiServer is the API server.
type ApiServer struct {
	server *server.Server
}

var skyflowConfigFilePath string = "./skyflow.yaml"

// NewApiServer creates a new API server.
func NewApiServer(configFilePath string) *ApiServer {
	if configFilePath != "" {
		trpc.ServerConfigPath = configFilePath
	}

	s := trpc.NewServer()
	pb.RegisterCommonServiceService(s, &CommonServiceHandler{})
	pb.RegisterSkyflowV1ServiceService(s, &SkyflowServiceHandler{})

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
