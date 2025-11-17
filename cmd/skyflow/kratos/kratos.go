package kratos

import (
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	id, _ = os.Hostname()
)

func MustNewApp(name string, version string, logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	id, err := os.Hostname()
	if err != nil {
		panic(err)
	}

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
