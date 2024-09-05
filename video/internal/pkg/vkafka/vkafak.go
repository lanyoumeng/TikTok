package vkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	userv1 "video/api/user/v1"
	"video/internal/conf"
	"video/internal/pkg/aliyun"
	"video/internal/pkg/model"
)

func NewVdKafkaSender(k *conf.Kafka) (*kafka.Writer, error) {
	topic := k.Video.Publish.Topic

	writer := &kafka.Writer{
		Addr:  kafka.TCP(k.Broker), //不定长参数，支持传入多个broker的ip:port
		Topic: topic,               //为所有message指定统一的topic。如果这里不指定统一的Topic，则创建kafka.Message{}时需要分别指定Topic

		Balancer:     &kafka.Hash{},     //把message的key进行hash，确定partition  &kafka.LeastBytes{}指定分区的balancer模式为最小字节分布  &kafka.RoundRobin{}：循环地将消息依次发送到每个分区，实现轮询效果。(默认)
		WriteTimeout: 1 * time.Second,   //设定写超时
		RequiredAcks: kafka.RequireNone, //RequireNone不需要等待ack返回，效率最高，安全性最低；
		// RequireOne只需要确保Leader写入成功就可以发送下一条消息；
		// kafka.RequireAll需要确保Leader和所有Follower都写入成功才可以发送下一条消息。
		AllowAutoTopicCreation: true, //Topic不存在时自动创建。生产环境中一般设为false，由运维管理员创建Topic并配置partition数目
		// Async:                  true, // 异步,在后台发送消息，而不会阻塞主线程。
		// Logger:      kafka.LoggerFunc(zap.NewExample().Sugar().Infof), //使用第三方日志库
		// ErrorLogger: kafka.LoggerFunc(zap.NewExample().Sugar().Errorf),
		// Compression: kafka.Snappy, //压缩
	}
	// defer writer.Close() //记得关闭连接

	return writer, nil
}

func NewVdKafkaReader(log *log.Helper, k *conf.Kafka, bucket *oss.Bucket, db *gorm.DB, rdb *redis.Client, uc userv1.UserServiceClient) *kafka.Reader {
	broker := k.Broker
	topic := k.Video.Publish.Topic
	groupId := k.Video.Publish.GroupId

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
	go InitVideoKafkaConsumer(log, reader, bucket, db, rdb, uc)

	return reader
}

// 初始化视频消息消费者
func InitVideoKafkaConsumer(log *log.Helper, reader *kafka.Reader, bucket *oss.Bucket, db *gorm.DB, rdb *redis.Client, uc userv1.UserServiceClient) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //注册信号2和15
		sig := <-c                                        //阻塞，直到信号的到来
		fmt.Printf("receive signal %s\n", sig.String())
		if reader != nil {
			err := reader.Close()
			if err != nil {
				return
			}
		}
		os.Exit(0) //进程退出
	}()

	ctx := context.Background()
	var cnt int
	for { //消息队列里随时可能有新消息进来，所以这里是死循环，类似于读Channel
		if message, err := reader.ReadMessage(ctx); err != nil {
			log.Debug("read message from kafka failed: %v", err)
			break
		} else {

			//1.kafka接收到视频信息后上传到阿里云oss
			fmt.Printf("topic=%s, partition=%d, offset=%d, key=%s, message content=%s\n", message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
			videoKafkaMessage := &model.VideoKafkaMessage{}
			if err := json.Unmarshal(message.Value, &videoKafkaMessage); err != nil {
				fmt.Printf("json.Unmarshal failed: %v", err)
			}

			cnt++
			log.Debug("vkafka/Message::::::::::::", videoKafkaMessage, " 		cnt:", cnt)
			if bucket == nil {
				log.Debugf("bucket is nil")
				continue
			}
			//log.Debug("bucket--------------", bucket)

			playURL, coverURL, err := aliyun.UploadFile(bucket, videoKafkaMessage)
			if err != nil {
				log.Debugf("upload file failed: %v", err)
			}

			//存储的video
			video := &model.Video{}
			video.AuthorId = videoKafkaMessage.AuthorId
			video.CoverUrl = coverURL
			video.PlayUrl = playURL
			video.Title = videoKafkaMessage.Title

			// 开始事务  分布式事务

			//2.执行数据库操作
			if err := db.Model(&model.Video{}).Save(&video).Error; err != nil {
				log.Info("Error creating user1:", err)
			}
			//3 Zset存储所有视频id 不过期
			currentTime := float64(video.CreatedAt.Unix())
			err = rdb.ZAdd(ctx, "videoAll", &redis.Z{Score: currentTime, Member: video.Id}).Err()

			//4.更新user服务 redis的work_count	字段
			userId := video.AuthorId
			var workCount int64
			if err := db.Model(&model.Video{}).Where("author_id = ?", userId).Count(&workCount).Error; err != nil {
				log.Error("Error creating user2:", err)
			}
			_, err = uc.UpdateWorkCnt(ctx, &userv1.UpdateWorkCntRequest{UserId: strconv.FormatInt(userId, 10), WorkCount: workCount})
			if err != nil {
				log.Error("Error creating user2:", err)
			}

			//// 4 再开启一个goroutine：将视频存到redis
			//key := "videoInfo::" + strconv.FormatInt(video.Id, 10)
			//videoMap, err := tool.StructToMap(video)
			//err = rdb.HSet(ctx, key, videoMap).Err()
			//if err != nil {
			//	log.Error("Error creating user2:", err)
			//}
			//_, err = rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Result()
			//if err != nil {
			//	log.Error("Error creating user2:", err)
			//}

			// 再开启一个goroutine：将视频id加入到布隆过滤器中

		}
	}

}
