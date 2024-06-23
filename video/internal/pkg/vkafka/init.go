package vkafka

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	userV1 "video/api/user/v1"
	"video/internal/conf"
)

// ProviderSet is vkafka providers.
var ProviderSet = wire.NewSet(InitKafkaConsumer)

func InitKafkaConsumer(logger log.Logger,
	k *conf.Kafka,
	db *gorm.DB,
	rdb *redis.Client,
	bucket *oss.Bucket,
	uc userV1.UserServiceClient) *kafka.Reader {
	log := log.NewHelper(log.With(logger, "module", "vkafka"))

	NewRedisKafkaReader(log, k, rdb)
	reader := NewVdKafkaReader(log, k, bucket, db, rdb, uc)

	return reader
}
