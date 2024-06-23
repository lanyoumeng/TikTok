package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v2"
)

// Config 结构体用于表示配置
type Config struct {
	Server struct {
		HTTP struct {
			Addr    string `yaml:"addr"`
			Timeout string `yaml:"timeout"`
		} `yaml:"http"`
		GRPC struct {
			Addr    string `yaml:"addr"`
			Timeout string `yaml:"timeout"`
		} `yaml:"grpc"`
	} `yaml:"server"`
	Service struct {
		User struct {
			Endpoint string `yaml:"endpoint"`
		} `yaml:"user"`
		Favorite struct {
			Endpoint string `yaml:"endpoint"`
		} `yaml:"favorite"`
		Video struct {
			Endpoint string `yaml:"endpoint"`
		} `yaml:"video"`
		Comment struct {
			Endpoint string `yaml:"endpoint"`
		} `yaml:"comment"`
		Relation struct {
			Endpoint string `yaml:"endpoint"`
		} `yaml:"relation"`
	} `yaml:"service"`
	Data struct {
		Database struct {
			Driver string `yaml:"driver"`
			Source string `yaml:"source"`
		} `yaml:"database"`
		Redis struct {
			Addr         string `yaml:"addr"`
			Password     string `yaml:"password"`
			DB           int    `yaml:"db"`
			DialTimeout  string `yaml:"dial_timeout"`
			ReadTimeout  string `yaml:"read_timeout"`
			WriteTimeout string `yaml:"write_timeout"`
		} `yaml:"redis"`
	} `yaml:"data"`
	Trace struct {
		Endpoint string `yaml:"endpoint"`
	} `yaml:"trace"`
	Auth struct {
		JWTKey string `yaml:"jwt_key"`
	} `yaml:"auth"`
	Etcd struct {
		Address string `yaml:"address"`
	} `yaml:"etcd"`
	Prometheus struct {
		Post string `yaml:"post"`
		Path string `yaml:"path"`
	} `yaml:"prometheus"`
	Kafka struct {
		Broker   string `yaml:"broker"`
		Consumer struct {
			Group string `yaml:"group"`
			Topic string `yaml:"topic"`
		} `yaml:"consumer"`
		Producer struct {
			Topic string `yaml:"topic"`
		} `yaml:"producer"`
	} `yaml:"kafka"`
	AliyunOSS struct {
		Endpoint        string `yaml:"endpoint"`
		AccessKeyID     string `yaml:"accessKeyId"`
		AccessKeySecret string `yaml:"accessKeySecret"`
		BucketName      string `yaml:"bucketName"`
	} `yaml:"aliyunOSS"`
}

func main() {
	// 读取配置文件
	configFile := "config.yaml"
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	// 连接etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.Etcd.Address},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer cli.Close()

	// 设置前缀
	prefix := "/TikTok-config/video/"

	// 存储配置到etcd
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	putConfig(cli, ctx, prefix+"server/http/addr", config.Server.HTTP.Addr)
	putConfig(cli, ctx, prefix+"server/http/timeout", config.Server.HTTP.Timeout)
	putConfig(cli, ctx, prefix+"server/grpc/addr", config.Server.GRPC.Addr)
	putConfig(cli, ctx, prefix+"server/grpc/timeout", config.Server.GRPC.Timeout)

	putConfig(cli, ctx, prefix+"service/user/endpoint", config.Service.User.Endpoint)
	putConfig(cli, ctx, prefix+"service/favorite/endpoint", config.Service.Favorite.Endpoint)
	putConfig(cli, ctx, prefix+"service/video/endpoint", config.Service.Video.Endpoint)
	putConfig(cli, ctx, prefix+"service/comment/endpoint", config.Service.Comment.Endpoint)
	putConfig(cli, ctx, prefix+"service/relation/endpoint", config.Service.Relation.Endpoint)

	putConfig(cli, ctx, prefix+"data/database/driver", config.Data.Database.Driver)
	putConfig(cli, ctx, prefix+"data/database/source", config.Data.Database.Source)
	putConfig(cli, ctx, prefix+"data/redis/addr", config.Data.Redis.Addr)
	putConfig(cli, ctx, prefix+"data/redis/password", config.Data.Redis.Password)
	putConfig(cli, ctx, prefix+"data/redis/db", fmt.Sprintf("%d", config.Data.Redis.DB))
	putConfig(cli, ctx, prefix+"data/redis/dial_timeout", config.Data.Redis.DialTimeout)
	putConfig(cli, ctx, prefix+"data/redis/read_timeout", config.Data.Redis.ReadTimeout)
	putConfig(cli, ctx, prefix+"data/redis/write_timeout", config.Data.Redis.WriteTimeout)

	putConfig(cli, ctx, prefix+"trace/endpoint", config.Trace.Endpoint)
	putConfig(cli, ctx, prefix+"auth/jwt_key", config.Auth.JWTKey)
	putConfig(cli, ctx, prefix+"prometheus/post", config.Prometheus.Post)
	putConfig(cli, ctx, prefix+"prometheus/path", config.Prometheus.Path)

	putConfig(cli, ctx, prefix+"kafka/broker", config.Kafka.Broker)
	putConfig(cli, ctx, prefix+"kafka/consumer/group", config.Kafka.Consumer.Group)
	putConfig(cli, ctx, prefix+"kafka/consumer/topic", config.Kafka.Consumer.Topic)
	putConfig(cli, ctx, prefix+"kafka/producer/topic", config.Kafka.Producer.Topic)

	putConfig(cli, ctx, prefix+"aliyunOSS/endpoint", config.AliyunOSS.Endpoint)
	putConfig(cli, ctx, prefix+"aliyunOSS/accessKeyId", config.AliyunOSS.AccessKeyID)
	putConfig(cli, ctx, prefix+"aliyunOSS/accessKeySecret", config.AliyunOSS.AccessKeySecret)
	putConfig(cli, ctx, prefix+"aliyunOSS/bucketName", config.AliyunOSS.BucketName)
}

func putConfig(cli *clientv3.Client, ctx context.Context, key, value string) {
	_, err := cli.Put(ctx, key, value)
	if err != nil {
		log.Fatalf("Failed to put key %s: %v", key, err)
	}
}
