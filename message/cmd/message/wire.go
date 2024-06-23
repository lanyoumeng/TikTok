//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"message/internal/biz"
	"message/internal/conf"
	"message/internal/data"
	"message/internal/pkg/mkafka"
	"message/internal/server"
	"message/internal/service"
)

// wireApp init kratos application.
func wireApp(
	*conf.Kafka,
	*conf.Prometheus,
	*conf.Etcd,
	*conf.Auth,
	*conf.Service,

	*conf.Server,
	*conf.Data,
	log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		mkafka.ProviderSet,
		newApp))
}
