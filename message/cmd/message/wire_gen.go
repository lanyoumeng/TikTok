// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"message/internal/biz"
	"message/internal/conf"
	"message/internal/data"
	"message/internal/pkg/mkafka"
	"message/internal/server"
	"message/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(kafka *conf.Kafka, prometheus *conf.Prometheus, etcd *conf.Etcd, auth *conf.Auth, confService *conf.Service, confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	registrar, cleanup, err := data.NewRegistrar(etcd, logger)
	if err != nil {
		return nil, nil, err
	}
	db := data.NewDB(confData)
	client := data.NewRedis(confData)
	discovery := data.NewDiscovery(etcd)
	favoriteServiceClient := data.NewFavoriteClient(discovery, confService)
	commentServiceClient := data.NewCommentClient(discovery, confService)
	relationServiceClient := data.NewRelationClient(discovery, confService)
	videoServiceClient := data.NewVideoClient(discovery, confService)
	reader := mkafka.InitKafkaConsumer(logger, kafka, client)
	dataData, cleanup2, err := data.NewData(logger, db, client, favoriteServiceClient, commentServiceClient, relationServiceClient, videoServiceClient, reader)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	messageRepo := data.NewMessageRepo(dataData, logger)
	messageUsecase := biz.NewMessageUsecase(messageRepo, logger)
	messageService := service.NewMessageService(messageUsecase, auth, logger)
	grpcServer := server.NewGRPCServer(prometheus, confServer, auth, messageService, logger)
	app := newApp(registrar, logger, grpcServer)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}