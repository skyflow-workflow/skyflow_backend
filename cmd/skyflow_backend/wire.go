//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/skyflow-workflow/skyflow_backbend/internal/biz"
	"github.com/skyflow-workflow/skyflow_backbend/internal/conf"
	"github.com/skyflow-workflow/skyflow_backbend/internal/data"
	"github.com/skyflow-workflow/skyflow_backbend/internal/server"
	"github.com/skyflow-workflow/skyflow_backbend/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
