package kratos

import (
	"log/slog"
	"os"

	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/skyflow-workflow/skyflow_backend/internal/conf"
	"github.com/skyflow-workflow/skyflow_backend/workflow"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

var (
	Name    string
	Version string
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	return srv
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	opts = append(opts, http.ResponseEncoder(CustomResponseEncoder))
	opts = append(opts, http.ErrorEncoder(CustomErrorHandler))
	srv := http.NewServer(opts...)
	return srv
}

func GetDefaultLogger() log.Logger {
	id := getHostname()
	defaultlogger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	return defaultlogger
}

func GetBootstrapConfig(configPath string) (*conf.Bootstrap, error) {
	c := config.New(
		config.WithSource(
			file.NewSource(configPath),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		return nil, err
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		return nil, err
	}
	return &bc, nil
}

// Server  kratos 服务
type Server struct {
	grpcServer *grpc.Server
	httpServer *http.Server
	wfSvc      workflow.WorkflowService
	logger     log.Logger
	app        *kratos.App
}

type ServerOption func(*Server)

type ServerOptions struct {
}

func NewServer(conf string, svc workflow.WorkflowService) (*Server, error) {

	logger := GetDefaultLogger()
	bootstrap, err := GetBootstrapConfig(conf)
	if err != nil {
		return nil, err
	}

	grpcServer := NewGRPCServer(bootstrap.Server)
	httpServer := NewHTTPServer(bootstrap.Server, logger)
	InitAppServer(grpcServer, httpServer, svc)

	// config slogger
	ConfigSlogger()

	app := NewApp(Name, Version, logger, grpcServer, httpServer)
	server := &Server{
		grpcServer: grpcServer,
		httpServer: httpServer,
		wfSvc:      svc,
		logger:     logger,
		app:        app,
	}
	return server, nil
}

func (s *Server) Start() error {
	return s.app.Run()
}
func (s *Server) Stop() error {
	return nil
}

func ConfigSlogger() {

	l := slog.New(
		slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			}),
	)
	slog.SetDefault(l)
}
