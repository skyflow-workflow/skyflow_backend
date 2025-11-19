package trpc

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"

	_ "trpc.group/trpc-go/trpc-gateway/plugin/accesslog"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/filter"
)

func TestInitTrpcConfig(t *testing.T) {
	configFilePath := "../mock/trpc_go.yaml"

	wd, err := os.Getwd()
	assert.Equal(t, err, nil)
	slog.Info("Current working directory:", "wd", wd)
	trpc.ServerConfigPath = configFilePath
	s := trpc.NewServer()
	slog.Info("TRPC server initialized",
		"server", s,
	)
	f := filter.GetServer("accesslog")
	if f == nil {
		t.Error("Expected accesslog filter to be registered, but it is nil")
	} else {
		slog.Info("Accesslog filter is registered",
			"filter", f,
		)
	}
	service := s.Service("trpc.http.skyflow")
	fmt.Println("Service Name:", service)
}
