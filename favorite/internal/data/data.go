package data

import (
	"context"
	"favorite/internal/conf"
	"favorite/internal/pkg/model"
	"time"

	// commentV1 "Favorite/api/comment/v1"
	// userV1 "Favorite/api/user/v1"
	userV1 "favorite/api/user/v1"
	videoV1 "favorite/api/video/v1"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/google/wire"
	"github.com/segmentio/kafka-go"

	grpcx "google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	NewRedis,
	NewBizFavoriteRepo,
	NewEtcdClient,

	NewDiscovery,
	NewRegistrar,
	NewVideoServiceClient,
	NewUserServiceClient,
)

// Data .
type Data struct {
	db  *gorm.DB
	rdb *redis.Client

	vc    videoV1.VideoServiceClient
	userc userV1.UserServiceClient

	LogKafkaReader *kafka.Reader
	// bucket        *oss.Bucket
}

// NewData .
func NewData(c *conf.Data,
	logger log.Logger,
	db *gorm.DB,
	rdb *redis.Client,

	userc userV1.UserServiceClient,
	vc videoV1.VideoServiceClient,

	LogKafkaReader *kafka.Reader,
) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
			db:    db,
			rdb:   rdb,
			userc: userc,
			vc:    vc,

			LogKafkaReader: LogKafkaReader,
			// bucket:        bucket,
		},
		cleanup,
		nil
}

// NewDB .
func NewDB(data *conf.Data) *gorm.DB {

	db, err := gorm.Open(mysql.Open(data.Database.Source), &gorm.Config{})
	if err != nil {

		panic("failed to connect database")
	}

	// 自动迁移模型，将模型的结构映射到数据库表中
	err = db.AutoMigrate(&model.Favorite{})
	if err != nil {
		log.Errorf("AutoMigrate model.Favorite failed: %v", err)
		return nil
	}

	return db
}

func NewRedis(conf *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password,
		DB:           int(conf.Redis.Db),
		DialTimeout:  conf.Redis.DialTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
	})
	rdb.AddHook(redisotel.TracingHook{})
	return rdb
}

func NewRegistrar(etcdpoint *conf.Etcd, logger log.Logger) (registry.Registrar, func(), error) {

	// ETCD源地址
	endpoint := []string{etcdpoint.Address}
	// ETCD配置信息
	etcdCfg := clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: 5 * time.Second,
		DialOptions: []grpcx.DialOption{grpcx.WithBlock()},
	}

	// 创建ETCD客户端
	client, err := clientv3.New(etcdCfg)
	if err != nil {
		panic(err)
	}
	clean := func() {
		_ = client.Close()
	}

	// 创建服务注册 reg
	regi := etcd.New(client)

	return regi, clean, nil
}

// 链接用户服务 grpc
func NewEtcdClient(etcdpoint *conf.Etcd) *clientv3.Client {
	// ETCD源地址
	endpoint := []string{etcdpoint.Address}
	// ETCD配置信息
	etcdCfg := clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: time.Second,
		DialOptions: []grpcx.DialOption{grpcx.WithBlock()},
	}

	// 创建ETCD客户端
	client, err := clientv3.New(etcdCfg)
	if err != nil {
		panic(err)
	}

	return client

}

// 链接用户服务 grpc
func NewDiscovery(etcdpoint *conf.Etcd) registry.Discovery {
	// ETCD源地址
	endpoint := []string{etcdpoint.Address}
	// ETCD配置信息
	etcdCfg := clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: time.Second,
		DialOptions: []grpcx.DialOption{grpcx.WithBlock()},
	}

	// 创建ETCD客户端
	client, err := clientv3.New(etcdCfg)
	if err != nil {
		panic(err)
	}

	// new dis with etcd client
	dis := etcd.New(client)
	return dis

}

func NewUserServiceClient(covery registry.Discovery, s *conf.Service) userV1.UserServiceClient {
	// ETCD源地址 discovery:///TikTok.user.service
	endpoint := s.User.Endpoint
	conn, err := grpc.DialInsecure( //不使用TLS
		context.Background(),
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(covery),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
		grpc.WithTimeout(2*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := userV1.NewUserServiceClient(conn)
	return c
}

func NewVideoServiceClient(covery registry.Discovery, s *conf.Service) videoV1.VideoServiceClient {
	// ETCD源地址 discovery:///TikTok.video.service
	endpoint := s.Video.Endpoint
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(covery),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
		grpc.WithTimeout(2*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),

		// grpc.WithTLSConfig(), // 使用TLS连接
	)

	if err != nil {
		panic(err)
	}
	c := videoV1.NewVideoServiceClient(conn)
	return c
}

// func NewCommentClient(covery registry.Discovery, s *conf.Service) commentV1.CommentServiceClient {
// 	// ETCD源地址 discovery:///TikTok.comment.service
// 	endpoint := s.Comment.Endpoint
// 	conn, err := grpc.Dial(
// 		context.Background(),
// 		grpc.WithEndpoint(endpoint),
// 		grpc.WithDiscovery(covery),
// 		grpc.WithMiddleware(
// 			tracing.Client(),
// 			recovery.Recovery(),
// 		),
// 		grpc.WithTimeout(2*time.Second),
// 		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	c := commentV1.NewCommentServiceClient(conn)
// 	return c
// }
