package trpc

import (
	"log/slog"

	"github.com/skyflow-workflow/skyflow_backend/pkg/trpclog"
	"github.com/skyflow-workflow/skyflow_backend/workflow"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/server"
)

type TrpcServer struct {
	trpc_conf string
	server    *server.Server
	wf        workflow.WorkflowService
}

func NewServer(trpc_conf string, wf workflow.WorkflowService) (*TrpcServer, error) {
	if trpc_conf != "" {
		trpc.ServerConfigPath = trpc_conf // Set the TRPC server configuration path
	}

	// load TRPC server configuration
	s := trpc.NewServer()

	// load trpc logger config and transform it to slog logger
	logger := trpclog.NewHandlerFromTrpcLogger(nil)
	// set slog as the default logger
	slog.SetDefault(slog.New(logger))

	slog.Info("Initializing TRPC server", "configPath", trpc_conf)

	InitTrpcServer(s, wf)

	ts := &TrpcServer{
		trpc_conf: trpc_conf,
		server:    s,
		wf:        wf,
	}
	return ts, nil
}

func InitTrpcServer(s *server.Server, wf workflow.WorkflowService) {
	// pbv1.RegisterCommonServiceService(server, &CommonServiceHandler{})
	// pbv1.RegisterSkyflowV1ServiceService(server, &SkyflowServiceHandler{
	// 	wfSvc: svc,
	// })
}

func (ts *TrpcServer) Start() error {
	// Start the API server
	return ts.server.Serve()
	// if err := ts.server.Serve(); err != nil {
	// 	log.Error(err)
	// }
}
func (ts *TrpcServer) Stop() error {
	return nil
}
