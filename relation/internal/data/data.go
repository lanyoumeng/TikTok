package data

import (
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
	commentV1 "relation/api/comment/v1"
	favoriteV1 "relation/api/favorite/v1"
	messageV1 "relation/api/message/v1"
	userV1 "relation/api/user/v1"
	videov1 "relation/api/video/v1"
	"relation/internal/conf"
	"relation/internal/pkg/model"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	NewRedis,

	NewRelationRepo,

	NewDiscovery,
	NewRegistrar,

	NewUserServiceClient,
	NewFavoriteClient,
	NewCommentClient,
	NewVideoClient,
	NewMessageClient,
)

// Data .
type Data struct {
	// TODO wrapped database client
	db  *gorm.DB
	rdb *redis.Client

	userc    userV1.UserServiceClient         // user client
	favc     favoriteV1.FavoriteServiceClient // favorite client
	commentc commentV1.CommentServiceClient   // comment client
	videoc   videov1.VideoServiceClient       // video client
	messagec messageV1.MessageServiceClient

	LogKafkaReader *kafka.Reader
}

// NewData .
func NewData(
	logger log.Logger,
	db *gorm.DB,
	rdb *redis.Client,

	userc userV1.UserServiceClient,
	favc favoriteV1.FavoriteServiceClient,
	commentc commentV1.CommentServiceClient,
	videoc videov1.VideoServiceClient,
	messagec messageV1.MessageServiceClient,

	LogKafkaReader *kafka.Reader,

) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db:  db,
		rdb: rdb,

		userc:    userc,
		favc:     favc,
		commentc: commentc,
		videoc:   videoc,
		messagec: messagec,

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
	err = db.AutoMigrate(&model.Relation{})
	if err != nil {
		log.Error(err)
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

	//if err := rdb.Close(); err != nil {
	//	log.Error(err)
	//}
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
		grpc.WithTimeout(2*time.Second),
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
		grpc.WithTimeout(2*time.Second),
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		panic(err)
	}
	c := favoriteV1.NewFavoriteServiceClient(conn)
	return c
}

func NewCommentClient(covery registry.Discovery, s *conf.Service) commentV1.CommentServiceClient {
	// ETCD源地址 discovery:///TikTok.comment.service
	endpoint := s.Comment.Endpoint
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
	)
	if err != nil {
		panic(err)
	}
	c := commentV1.NewCommentServiceClient(conn)
	return c
}

func NewVideoClient(covery registry.Discovery, s *conf.Service) videov1.VideoServiceClient {
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
	)
	if err != nil {
		panic(err)
	}
	c := videov1.NewVideoServiceClient(conn)
	return c
}

func NewMessageClient(covery registry.Discovery, s *conf.Service) messageV1.MessageServiceClient {
	// ETCD源地址 discovery:///TikTok.message.service
	endpoint := s.Message.Endpoint
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
	)
	if err != nil {
		panic(err)
	}
	c := messageV1.NewMessageServiceClient(conn)
	return c
}
