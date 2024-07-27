package rkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"relation/internal/conf"
	"relation/internal/pkg/model"
	"strconv"
	"syscall"
	"time"
)

type RedisKafkaMessage struct {
	Before model.Relation `json:"before"`
	After  model.Relation `json:"after"`
	Op     string         `json:"op"`
}

func NewRedisKafkaReader(log *log.Helper, k *conf.Kafka, rdb *redis.Client) *kafka.Reader {
	broker := k.Broker
	topic := k.User.User.Topic
	groupId := k.User.User.GroupId

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker}, //支持传入多个broker的ip:port
		Topic:   topic,
		// Partition: 1,  //注意partition和groupID不能同时设置

		GroupID: groupId, //一个Group内消费到的消息不会重复
		// 在使用消费者组时会有以下限制：
		// - `(*Reader).SetOffset` 当设置了GroupID时会返回错误
		// - `(*Reader).Offset` 当设置了GroupID时会永远返回 `-1`
		// - `(*Reader).Lag` 当设置了GroupID时会永远返回 `-1`
		// - `(*Reader).ReadLag` 当设置了GroupID时会返回错误
		// - `(*Reader).Stats` 当设置了GroupID时会返回一个`-1`的分区
		CommitInterval: 1 * time.Second,   //每隔多长时间自动commit一次offset。即一边读一边向kafka上报读到了哪个位置。
		StartOffset:    kafka.FirstOffset, //当一个特定的partition没有commited offset时(比如第一次读一个partition，之前没有commit过)，通过StartOffset指定从第一个还是最后一个位置开始消费。StartOffset的取值要么是FirstOffset要么是LastOffset，LastOffset表示Consumer启动之前生成的老数据不管了。仅当指定了GroupID时，StartOffset才生效。
		// MaxBytes:       10e6,              // 10MB
		// 	Logger:      kafka.LoggerFunc(logf), //以自定义一个Logger或使用第三方日志库
		// ErrorLogger: kafka.LoggerFunc(logf),
	})

	//协程启动消费者
	go InitRedisKafkaConsumer(context.Background(), log, reader, rdb)
	return reader
}

// 初始化缓存消息消费者
func InitRedisKafkaConsumer(ctx context.Context, log *log.Helper, reader *kafka.Reader, rdb *redis.Client) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //注册信号2和15
		sig := <-c                                        //阻塞，直到信号的到来
		fmt.Printf("receive signal %s\n", sig.String())
		if reader != nil {
			reader.Close()
		}
		os.Exit(0) //进程退出
	}()

	for { //消息队列里随时可能有新消息进来，所以这里是死循环，类似于读Channel
		if message, err := reader.ReadMessage(ctx); err != nil {
			log.Errorf("read message from kafka failed: %v", err)
			break
		} else {

			// fmt.Printf("topic=%s, partition=%d, offset=%d, key=%s, message content=%s\n", message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
			redisKafkaMessage := RedisKafkaMessage{}
			if err := json.Unmarshal(message.Value, &redisKafkaMessage); err != nil {
				log.Errorf("json.Unmarshal failed: %v", err)
			}

			//关注用户的id  set  follow::user_id 关注用户的ids
			//粉丝的id  set  follower::user_id  粉丝的ids

			//relation表创建一条记录：+过期时间
			//添加对应key-member
			//
			//relation表不会更新
			//
			//relation表删除一条记录： +过期时间
			//移除对应key-member

			if redisKafkaMessage.Op == "c" {

				//关注用户的id  set  follow::user_id 关注用户的ids
				//粉丝的id  set  follower::user_id  粉丝的ids
				userId := redisKafkaMessage.After.UserId
				toUserId := redisKafkaMessage.After.ToUserId

				err := rdb.SAdd(ctx, "follow::"+strconv.FormatInt(userId, 10), redisKafkaMessage.After.ToUserId).Err()
				if err != nil {
					log.Errorf("SAdd failed: %v", err)
				}
				//过期时间
				err = rdb.Expire(ctx, "follow::"+strconv.FormatInt(userId, 10), 24*time.Hour).Err()
				if err != nil {
					log.Errorf("Expire failed: %v", err)
				}
				err = rdb.SAdd(ctx, "follower::"+strconv.FormatInt(toUserId, 10), redisKafkaMessage.After.UserId).Err()
				if err != nil {
					log.Errorf("SAdd failed: %v", err)
				}
				//过期时间
				err = rdb.Expire(ctx, "follower::"+strconv.FormatInt(toUserId, 10), 24*time.Hour).Err()
				if err != nil {
					log.Errorf("Expire failed: %v", err)
				}
			}

			if redisKafkaMessage.Op == "d" {
				userId := redisKafkaMessage.Before.UserId
				toUserId := redisKafkaMessage.Before.ToUserId

				err := rdb.SRem(ctx, "follow::"+strconv.FormatInt(userId, 10), redisKafkaMessage.Before.ToUserId).Err()
				if err != nil {
					log.Errorf("SRem failed: %v", err)
				}
				//过期时间
				err = rdb.Expire(ctx, "follow::"+strconv.FormatInt(userId, 10), 24*time.Hour).Err()
				if err != nil {
					log.Errorf("Expire failed: %v", err)
				}
				err = rdb.SRem(ctx, "follower::"+strconv.FormatInt(toUserId, 10), redisKafkaMessage.Before.UserId).Err()
				if err != nil {
					log.Errorf("SRem failed: %v", err)
				}
				//过期时间
				err = rdb.Expire(ctx, "follower::"+strconv.FormatInt(toUserId, 10), 24*time.Hour).Err()
				if err != nil {
					log.Errorf("Expire failed: %v", err)
				}
			}

		}
	}

}
