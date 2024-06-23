package favkafka

import (
	videoV1 "favorite/api/video/v1"
	"favorite/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/segmentio/kafka-go"
)

// ProviderSet is vkafka providers.
var ProviderSet = wire.NewSet(InitKafkaConsumer)

func InitKafkaConsumer(
	logger log.Logger,
	k *conf.Kafka,
	rdb *redis.Client,
	vc videoV1.VideoServiceClient,
) *kafka.Reader {

	log := log.NewHelper(log.With(logger, "module", "vkafka"))

	reader := NewRedisKafkaReader(log, k, rdb, vc)

	return reader
}
