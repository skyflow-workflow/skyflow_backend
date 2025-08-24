package cmd

import (
	"fmt"
	"log/slog"
	"testing"

	_ "trpc.group/trpc-go/trpc-gateway/plugin/accesslog"
	"trpc.group/trpc-go/trpc-go/filter"
)

func TestInitTrpcConfig(t *testing.T) {
	configFilePath := "./mock/configfile/trpc_go.yaml"

	s := InitializeTrpcSever(configFilePath)
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
