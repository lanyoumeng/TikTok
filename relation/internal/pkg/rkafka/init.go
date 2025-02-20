package rkafka

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/segmentio/kafka-go"
	"relation/internal/conf"
)

// ProviderSet is vkafka providers.
var ProviderSet = wire.NewSet(InitKafkaConsumer)

func InitKafkaConsumer(
	logger log.Logger,
	k *conf.Kafka,
	rdb *redis.Client,
) *kafka.Reader {

	//log := log.NewHelper(log.With(logger, "module", "rkafka"))
	log := log.NewHelper(logger)
	reader := NewRedisKafkaReader(log, k, rdb)

	return reader
}
