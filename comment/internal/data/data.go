package data

import (
	favoriteV1 "comment/api/favorite/v1"
	relationV1 "comment/api/relation/v1"
	userV1 "comment/api/user/v1"
	videoV1 "comment/api/video/v1"

	"comment/internal/conf"
	"comment/internal/pkg/model"
	"context"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/segmentio/kafka-go"
	clientv3 "go.etcd.io/etcd/client/v3"
	grpcx "google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewCommentRepo,

	NewData,
	NewDB,
	NewRedis,

	NewDiscovery,
	NewRegistrar,

	NewUserServiceClient,
	NewFavoriteClient,
	NewRelationClient,
	NewVideoClient,
)

// Data .
type Data struct {
	// TODO wrapped database client
	db  *gorm.DB
	rdb *redis.Client

	userc     userV1.UserServiceClient         // user client
	favc      favoriteV1.FavoriteServiceClient // favorite client
	relationc relationV1.RelationServiceClient // relation client
	videoc    videoV1.VideoServiceClient       // video client

	//NewLogKafkaWriter   *kafka.Writer
	LogKafkaReader *kafka.Reader
}

// NewData .
func NewData(
	logger log.Logger,
	db *gorm.DB,
	rdb *redis.Client,

	favc favoriteV1.FavoriteServiceClient,
	userc userV1.UserServiceClient,
	relationc relationV1.RelationServiceClient,
	videoc videoV1.VideoServiceClient,

	LogKafkaReader *kafka.Reader,

) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db:  db,
		rdb: rdb,

		favc:      favc,
		userc:     userc,
		relationc: relationc,
		videoc:    videoc,

		LogKafkaReader: LogKafkaReader,
	}, cleanup, nil
}

// NewDB .
func NewDB(conf *conf.Data) *gorm.DB {

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, //禁用为关联创建外键约束
	})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移模型，将模型的结构映射到数据库表中
	err = db.AutoMigrate(&model.Comment{})
	if err != nil {
		log.Errorf("AutoMigrate model.Comment failed: %v", err)
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

//// redis随机过期时间
//func GetRandomExpireTime() time.Duration {
//	return time.Duration(300+rand.Intn(300)) * time.Second
//}

func NewRegistrar(etcdpoint *conf.Etcd, logger log.Logger) (registry.Registrar, func(), error) {
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
	clean := func() {
		_ = client.Close()
	}

	// 创建服务注册 reg
	regi := etcd.New(client)

	return regi, clean, nil
}

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
		grpc.WithTimeout(10*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := userV1.NewUserServiceClient(conn)
	return c
}

func NewFavoriteClient(covery registry.Discovery, s *conf.Service) favoriteV1.FavoriteServiceClient {
	// ETCD源地址 discovery:///TikTok.favorite.service
	endpoint := s.Favorite.Endpoint
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(covery),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
		grpc.WithTimeout(10*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := favoriteV1.NewFavoriteServiceClient(conn)
	return c
}

func NewRelationClient(covery registry.Discovery, s *conf.Service) relationV1.RelationServiceClient {
	// ETCD源地址 discovery:///TikTok.relation.service
	endpoint := s.Relation.Endpoint
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(covery),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery(),
		),
		grpc.WithTimeout(10*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := relationV1.NewRelationServiceClient(conn)
	return c

}

func NewVideoClient(covery registry.Discovery, s *conf.Service) videoV1.VideoServiceClient {
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
		grpc.WithTimeout(10*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := videoV1.NewVideoServiceClient(conn)
	return c
}
