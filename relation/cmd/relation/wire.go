//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"relation/internal/biz"
	"relation/internal/conf"
	"relation/internal/data"
	"relation/internal/pkg/rkafka"
	"relation/internal/server"
	"relation/internal/service"
)

// wireApp init kratos application.
func wireApp(
	*conf.Kafka,
	*conf.Prometheus,
	*conf.Etcd,
	*conf.Auth,
	*conf.Service, *conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {

	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,

		//因为Wire 期望所有的提供者集合都会在依赖图中被使用，
		// 但InitKafkaConsumer函数并没有返回值调用，
		//所以在Data结构体添加了VideoKafkaConsumer *kafka.Reader（init添加了*kafka.Reader返回值），
		//来让InitKafkaConsumer可以被添加到wire中
		rkafka.ProviderSet,
		newApp))
}
