package kratos

import (
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	v1 "github.com/skyflow-workflow/skyflow_backbend/api/v1"
	"github.com/skyflow-workflow/skyflow_backbend/server/apiserver"
	"github.com/skyflow-workflow/skyflow_backbend/workflow"

	"github.com/go-kratos/kratos/v2/log"
)

// getHostname 安全地获取主机名，避免panic
func getHostname() string {
	if hostname, err := os.Hostname(); err != nil {
		return "unknown"
	} else {
		return hostname
	}
}

func NewApp(name string, version string, logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	id := getHostname()

	return kratos.New(
		kratos.ID(id),
		kratos.Name(name),
		kratos.Version(version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func InitAppServer(gs *grpc.Server, hs *http.Server, wfSvc workflow.WorkflowService) {

	commonhandler := &apiserver.CommonServiceHandler{}
	skyflowhandler := apiserver.NewSkyflowServiceHandler(wfSvc)
	v1.RegisterCommonServiceServer(gs, commonhandler)
	v1.RegisterSkyflowV1ServiceServer(gs, skyflowhandler)
	v1.RegisterSkyflowV1ServiceHTTPServer(hs, skyflowhandler)
}
