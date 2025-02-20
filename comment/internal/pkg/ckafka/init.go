package ckafka

import (
	"comment/internal/conf"
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
) *kafka.Reader {

	//log := log.NewHelper(log.With(logger, "module", "ckafka"))
	log := log.NewHelper(logger)
	reader := NewRedisKafkaReader(log, k, rdb)

	return reader
}
