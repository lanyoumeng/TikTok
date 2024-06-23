package data

import (
	"favorite/internal/conf"
	"time"

	"github.com/segmentio/kafka-go"
)

func NewVdKafkaSender(k *conf.Kafka) (*kafka.Writer, error) {
	topic := k.Producer.Topic

	writer := &kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"), //不定长参数，支持传入多个broker的ip:port
		Topic: topic,                       //为所有message指定统一的topic。如果这里不指定统一的Topic，则创建kafka.Message{}时需要分别指定Topic

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
	defer writer.Close() //记得关闭连接

	return writer, nil
}

func NewVdKafkaReader(k *conf.Kafka) (*kafka.Reader, error) {
	broker := k.Broker
	topic := k.Consumer.Topic
	groupId := k.Consumer.Group

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
	// var FavoriteRepo FavoriteRepo
	// go FavoriteRepo.InitfavoriteKafkaConsumer(context.Background(), reader)
	return reader, nil
}

// // 初始化视频消息消费者
// func (v *FavoriteRepo) InitfavoriteKafkaConsumer(ctx context.Context, reader *kafka.Reader) {
// 	go func() {
// 		c := make(chan os.Signal, 1)
// 		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //注册信号2和15
// 		sig := <-c                                        //阻塞，直到信号的到来
// 		fmt.Printf("receive signal %s\n", sig.String())
// 		if reader != nil {
// 			reader.Close()
// 		}
// 		os.Exit(0) //进程退出
// 	}()

// 	for { //消息队列里随时可能有新消息进来，所以这里是死循环，类似于读Channel
// 		if message, err := reader.ReadMessage(ctx); err != nil {

// 			fmt.Printf("read message from kafka failed: %v", err)
// 			break
// 		} else {
// 			//kafka接收到视频信息后上传到阿里云oss

// 			// fmt.Printf("topic=%s, partition=%d, offset=%d, key=%s, message content=%s\n", message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
// 			var favoriteKafkaMessage *model.favoriteKafkaMessage
// 			if err := sonic.Unmarshal(message.Value, favoriteKafkaMessage); err != nil {
// 				fmt.Printf("sonic.Unmarshal failed: %v", err)
// 			}

// 			// playURL, coverURL, err := v.UploadFile(favoriteKafkaMessage)
// 			playURL, coverURL, err := UploadFile(v.data.bucket, favoriteKafkaMessage)
// 			if err != nil {
// 				fmt.Printf("UploadFile failed: %v", err)
// 			}

// 			//存储的favorite
// 			var favorite *model.favorite
// 			favorite.AuthorId = favoriteKafkaMessage.AuthorId
// 			favorite.CoverUrl = coverURL
// 			favorite.PlayUrl = playURL
// 			favorite.Title = favoriteKafkaMessage.Title

// 			if err = v.Save(ctx, favorite); err != nil {
// 				fmt.Printf("Save failed: %v", err)
// 			}

// 			// 	再开启一个goroutine：将视频上传到redis
// 			// 再开启一个goroutine：删除用户哈希字段 比如用户作品总数
// 			// 再开启一个goroutine：将视频id加入到布隆过滤器中
// 		}
// 	}

// }
