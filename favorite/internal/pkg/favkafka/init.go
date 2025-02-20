package favkafka

import (
	videoV1 "favorite/api/video/v1"
	"favorite/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/segmentio/kafka-go"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// ProviderSet is vkafka providers.
var ProviderSet = wire.NewSet(InitKafkaConsumer)

func InitKafkaConsumer(
	logger log.Logger,
	k *conf.Kafka,
	rdb *redis.Client,
	client *clientv3.Client,
	vc videoV1.VideoServiceClient,
) *kafka.Reader {

	//log := log.NewHelper(log.With(logger, "module", "vkafka"))
	log := log.NewHelper(logger)
	reader := NewRedisKafkaReader(log, k, client, rdb, vc)

	return reader
}
