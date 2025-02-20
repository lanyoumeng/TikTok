package ckafka

import (
	"comment/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
	"time"
)

// KafkaWriter 是一个将日志发送到 Kafka 的 io.Writer 实现
type KafkaWriter struct {
	Writer *kafka.Writer
}

// Write 实现 io.Writer 接口，将日志消息写入 Kafka
func (kw *KafkaWriter) Write(p []byte) (n int, err error) {

	message := kafka.Message{
		Value: p,
	}
	//log.Printf("Attempting to write message: %s", string(p))

	for i := 0; i < 3; i++ { //允许重试3次
		if err := kw.Writer.WriteMessages(context.Background(), //批量写入消息，原子操作，要么全写成`功，要么全写失败
			message,
		); err != nil {
			// if err == kafka.LeaderNotAvailable || errors.Is(err, context.DeadlineExceeded) {
			if err == kafka.LeaderNotAvailable { //首次写一个新的Topic时，会发生LeaderNotAvailable错误，重试一次就好了
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				log.Errorf("batch write message failed: %v", err)
			}
		} else {
			break //只要成功一次就不再尝试下一次了
		}
	}

	//log.Printf("Successfully wrote message: %s", string(p))
	return len(p), nil
}

// NewLogKafkaWriter 创建并返回一个新的 KafkaWriter 实例
func NewLogKafkaWriter(k *conf.Kafka) *KafkaWriter {

	topic := k.Log.Topic

	writer := &kafka.Writer{
		Addr:  kafka.TCP(k.Broker), //不定长参数，支持传入多个broker的ip:port
		Topic: topic,               //为所有message指定统一的topic。如果这里不指定统一的Topic，则创建kafka.Message{}时需要分别指定Topic

		//Balancer:     &kafka.Hash{},     //把message的key进行hash，确定partition  &kafka.LeastBytes{}指定分区的balancer模式为最小字节分布  &kafka.RoundRobin{}：循环地将消息依次发送到每个分区，实现轮询效果。(默认)
		WriteTimeout: 1 * time.Second, //设定写超时
		RequiredAcks: kafka.RequireNone,
		//RequireNone不需要等待ack返回，效率最高，安全性最低；
		// RequireOne只需要确保Leader写入成功就可以发送下一条消息；
		// kafka.RequireAll需要确保Leader和所有Follower都写入成功才可以发送下一条消息。
		AllowAutoTopicCreation: true, //Topic不存在时自动创建。生产环境中一般设为false，由运维管理员创建Topic并配置partition数目
		Async:                  true, // 异步,在后台发送消息，而不会阻塞主线程。
		// Logger:      kafka.LoggerFunc(zap.NewExample().Sugar().Infof), //使用第三方日志库
		// ErrorLogger: kafka.LoggerFunc(zap.NewExample().Sugar().Errorf),
		Compression: kafka.Snappy, //压缩
	}
	// defer writer.Close() //记得关闭连接

	return &KafkaWriter{
		Writer: writer,
	}
}
