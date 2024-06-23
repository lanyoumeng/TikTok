//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"favorite/internal/biz"
	"favorite/internal/conf"
	"favorite/internal/data"
	"favorite/internal/pkg/favkafka"
	"favorite/internal/server"
	"favorite/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Kafka, *conf.Service, *conf.Prometheus, *conf.Etcd, *conf.Auth, *conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		favkafka.ProviderSet,
		service.ProviderSet,
		newApp))
}
