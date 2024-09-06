package favkafka

import (
	"context"
	"encoding/json"
	videoV1 "favorite/api/video/v1"
	"favorite/internal/conf"
	"favorite/internal/pkg/model"
	"favorite/pkg/tool"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RedisKafkaMessage struct {
	Before model.Favorite `json:"before"`
	After  model.Favorite `json:"after"`
	Op     string         `json:"op"`
}

func NewRedisKafkaReader(log *log.Helper, k *conf.Kafka, client *clientv3.Client, rdb *redis.Client, vc videoV1.VideoServiceClient) *kafka.Reader {
	broker := k.Broker
	topic := k.Favorite.Favorite.Topic
	groupId := k.Favorite.Favorite.GroupId

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
	go InitRedisKafkaConsumer(context.Background(), log, reader, client, rdb, vc)
	return reader
}

// 初始化缓存消息消费者
func InitRedisKafkaConsumer(ctx context.Context, log *log.Helper, reader *kafka.Reader, etcdClient *clientv3.Client, rdb *redis.Client, vc videoV1.VideoServiceClient) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //注册信号2和15
		sig := <-c                                        //阻塞，直到信号的到来
		log.Debug("receive signal %s\n", sig.String())
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
			videoId := redisKafkaMessage.Before.VideoId
			userId := redisKafkaMessage.Before.UserId

			//etcd 创建租约并获取锁
			leaseResp, err := etcdClient.Grant(ctx, 10) // 租约 10 秒
			if err != nil {
				log.Fatalf("无法创建租约: %v", err)
			}

			lockName := fmt.Sprintf("lock_%d", videoId)
			key := fmt.Sprintf("locks/%s", lockName)

			_, err = etcdClient.Put(ctx, key, "locked", clientv3.WithLease(leaseResp.ID))
			if err != nil {
				log.Fatalf("无法获取锁: %v", err)
			}
			fmt.Println("成功获取锁")

			//favorite表创建/更新记录，
			if redisKafkaMessage.Op == "u" || redisKafkaMessage.Op == "c" {
				//点赞操作
				if redisKafkaMessage.After.Liked == true {

					g, ctx := errgroup.WithContext(ctx)

					g.Go(func() error {
						//1.string fav::<videoId>::<userId> true/false
						key := fmt.Sprintf("fav::%d::%d", videoId, userId)
						if err := rdb.Set(ctx, key, "true", tool.GetRandomExpireTime()).Err(); err != nil {
							log.Errorf("set redis failed: %v", err)
							return err
						}
						return nil
					})
					//2.set userFavList::<userId> video_id
					g.Go(func() error {
						key := fmt.Sprintf("userFavList::%d", userId)
						if err := rdb.SAdd(ctx, key, videoId).Err(); err != nil {
							log.Errorf("set redis failed: %v", err)
							return err
						}
						//过期时间
						if err := rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Err(); err != nil {
							log.Errorf("set redis failed: %v", err)
							return err
						}
						return nil
					})
					g.Go(func() error {
						//3.string  videoFavCnt::<videoId> int64
						key := fmt.Sprintf("videoFavCnt::%d", videoId)
						if err := rdb.Incr(ctx, key).Err(); err != nil {
							log.Errorf("set redis failed: %v", err)
							return err
						}
						if err := rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Err(); err != nil {
							log.Errorf("set redis failed: %v", err)
							return err
						}
						return nil
					})
					g.Go(func() error {
						//4.string	userTotalFavorited::userId int64
						//视频作者的获赞数+1

						//获取作者id
						resp, err := vc.GetAIdByVId(ctx, &videoV1.GetAIdByVIdReq{VideoId: videoId})
						if err != nil {
							return err
							log.Errorf("set redis failed: %v", err)
						}
						authorId := resp.AuthorId

						key := fmt.Sprintf("userTotalFavorited::%d", authorId)
						if err := rdb.Incr(ctx, key).Err(); err != nil {
							return err
							log.Errorf("set redis failed: %v", err)
						}
						if err := rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Err(); err != nil {
							log.Errorf("set redis failed: %v", err)
							return err
						}
						return nil
					})

					if err := g.Wait(); err != nil {
						log.Debug("set redis failed: %v", err)

					}

				}

				//取消点赞操作
				if redisKafkaMessage.After.Liked == false {
					g, ctx := errgroup.WithContext(ctx)

					g.Go(func() error {
						//1.string fav::<videoId>::<userId> true/false
						key := fmt.Sprintf("fav::%d::%d", videoId, userId)
						if err := rdb.Set(ctx, key, "false", tool.GetRandomExpireTime()).Err(); err != nil {
							log.Errorf("set redis failed: %v", err)
							return err
						}
						return nil
					})
					//2.set userFavList::<userId> video_id
					g.Go(func() error {
						key := fmt.Sprintf("userFavList::%d", userId)
						if err := rdb.SRem(ctx, key, videoId).Err(); err != nil {
							log.Errorf("set redis failed: %v", err)
							return err
						}
						//过期时间
						if err := rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Err(); err != nil {
							log.Errorf("set redis failed: %v", err)
							return err
						}
						return nil
					})
					g.Go(func() error {
						//3.string  videoFavCnt::<videoId> int64
						key := fmt.Sprintf("videoFavCnt::%d", videoId)
						if err := rdb.Decr(ctx, key).Err(); err != nil {
							log.Errorf("set redis failed: %v", err)
							return err
						}
						if err := rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Err(); err != nil {
							log.Errorf("set redis failed: %v", err)
							return err
						}
						return nil
					})
					g.Go(func() error {
						//4.string	userTotalFavorited::userId int64
						//获取作者id
						resp, err := vc.GetAIdByVId(ctx, &videoV1.GetAIdByVIdReq{VideoId: videoId})
						if err != nil {

							log.Errorf("set redis failed: %v", err)
							return err
						}
						authorId := resp.AuthorId
						key := fmt.Sprintf("userTotalFavorited::%d", authorId)
						if err := rdb.Decr(ctx, key).Err(); err != nil {

							log.Errorf("set redis failed: %v", err)
							return err
						}
						if err := rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Err(); err != nil {
							log.Errorf("set redis failed: %v", err)
							return err
						}
						return nil
					})

					if err := g.Wait(); err != nil {
						log.Errorf("set redis failed: %v", err)
					}

				}

			}

			//favorite表删除记录,点赞为true
			//数值-1
			if redisKafkaMessage.Op == "d" && redisKafkaMessage.Before.Liked == true {

				g, ctx := errgroup.WithContext(ctx)

				g.Go(func() error {
					//1.string fav::<videoId>::<userId> true/false
					key := fmt.Sprintf("fav::%d::%d", videoId, userId)
					if err := rdb.Del(ctx, key).Err(); err != nil {
						log.Errorf("set redis failed: %v", err)
						return err
					}
					return nil
				})
				//2.set userFavList::<userId> video_id
				g.Go(func() error {
					key := fmt.Sprintf("userFavList::%d", userId)
					if err := rdb.SRem(ctx, key, videoId).Err(); err != nil {
						log.Errorf("set redis failed: %v", err)
						return err
					}
					//过期时间
					if err := rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Err(); err != nil {
						log.Errorf("set redis failed: %v", err)
						return err
					}
					return nil
				})
				g.Go(func() error {
					//3.string  videoFavCnt::<videoId> int64
					key := fmt.Sprintf("videoFavCnt::%d", videoId)
					if err := rdb.Decr(ctx, key).Err(); err != nil {
						log.Errorf("set redis failed: %v", err)
						return err
					}
					if err := rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Err(); err != nil {
						log.Errorf("set redis failed: %v", err)
						return err
					}
					return nil
				})
				g.Go(func() error {
					//4.string	userTotalFavorited::userId int64
					key := fmt.Sprintf("userTotalFavorited::%d", userId)
					if err := rdb.Decr(ctx, key).Err(); err != nil {
						return err
						log.Errorf("set redis failed: %v", err)
					}
					if err := rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Err(); err != nil {
						log.Errorf("set redis failed: %v", err)
						return err
					}
					return nil
				})

				if err := g.Wait(); err != nil {
					log.Errorf("set redis failed: %v", err)
				}

			}

			//释放锁
			_, err = etcdClient.Delete(ctx, key)
			if err != nil {
				log.Fatalf("无法释放锁: %v", err)
			}

		}
	}

}
