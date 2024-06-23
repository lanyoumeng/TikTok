//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"comment/internal/biz"
	"comment/internal/conf"
	"comment/internal/data"
	"comment/internal/pkg/ckafka"
	"comment/internal/server"
	"comment/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Prometheus, *conf.Kafka, *conf.Etcd, *conf.Auth, *conf.Server, *conf.Service, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		ckafka.ProviderSet,
		newApp))
}
