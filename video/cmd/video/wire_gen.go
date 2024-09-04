// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"video/internal/biz"
	"video/internal/conf"
	"video/internal/data"
	"video/internal/pkg/aliyun"
	"video/internal/pkg/vkafka"
	"video/internal/server"
	"video/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(prometheus *conf.Prometheus, aliyunOSS *conf.AliyunOSS, kafka *conf.Kafka, etcd *conf.Etcd, auth *conf.Auth, confServer *conf.Server, confService *conf.Service, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	registrar, cleanup, err := data.NewRegistrar(etcd, logger)
	if err != nil {
		return nil, nil, err
	}
	db := data.NewDB(confData)
	client := data.NewRedis(confData)
	discovery := data.NewDiscovery(etcd)
	userServiceClient := data.NewUserServiceClient(discovery, confService)
	favoriteServiceClient := data.NewFavoriteClient(discovery, confService)
	commentServiceClient := data.NewCommentClient(discovery, confService)
	relationServiceClient := data.NewRelationClient(discovery, confService)
	writer, err := vkafka.NewVdKafkaSender(kafka)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	bucket, err := aliyun.NewBucket(aliyunOSS)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	reader := vkafka.InitKafkaConsumer(logger, kafka, db, client, bucket, userServiceClient)
	dataData, cleanup2, err := data.NewData(logger, db, client, userServiceClient, favoriteServiceClient, commentServiceClient, relationServiceClient, writer, reader, bucket)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	videoRepo := data.NewBizVideoRepo(dataData, logger)
	videoUsecase := biz.NewVideoUsecase(videoRepo, logger)
	videoService := service.NewVideoService(videoUsecase, auth, logger)
	grpcServer := server.NewGRPCServer(prometheus, confServer, auth, videoService, logger)
	httpServer := server.NewHTTPServer(confServer, videoService, auth, logger)
	app := newApp(registrar, logger, grpcServer, httpServer)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
