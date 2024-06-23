//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"user/internal/biz"
	"user/internal/conf"
	"user/internal/data"
	"user/internal/pkg/ukafka"
	"user/internal/server"
	"user/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Kafka, *conf.Prometheus, *conf.Etcd, *conf.Auth, *conf.Service, *conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		ukafka.ProviderSet,
		newApp))
}
