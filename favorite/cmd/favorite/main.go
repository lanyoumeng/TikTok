package main

import (
	"favorite/internal/conf"
	"favorite/internal/pkg/favkafka"
	"flag"
	knacos "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"os"

	"go.opentelemetry.io/otel/exporters/jaeger"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "TikTok.favorite.service"
	// Version is the version of the compiled software.
	Version = "v1"
	// flagconf is the config flag.
	flagconf string

	id = idFunc()
)

func idFunc() string {
	hostname, _ := os.Hostname()
	id := hostname + Name
	return id
}

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

// 设置全局trace
func initTracer(url string) error {
	// 创建 Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}

	tp := tracesdk.NewTracerProvider(
		// 将基于父span的采样率设置为100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// 始终确保在生产中批量处理
		tracesdk.WithBatcher(exp),
		// 在资源中记录有关此应用程序的信息
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(Name),
			attribute.String("env", "dev"),
			attribute.String("exporter", "jaeger"),
			attribute.Float64("float", 312.23),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

func newApp(r registry.Registrar, logger log.Logger, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
		),
		kratos.Registrar(r),
	)
}

func main() {
	flag.Parse()

	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         "public", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "../../deply/nacos/logs",
		CacheDir:            "../../deply/nacos/cache",
		LogLevel:            "debug",
	}

	// a more graceful way to create naming client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Debug(err)
	}

	c := config.New(
		config.WithSource(
			knacos.NewConfigSource(
				client,
				knacos.WithGroup("TikTok"),
				knacos.WithDataID("favorite.yaml"),
			),
			//file.NewSource(flagconf),
			// source,
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		log.Debug("c.Load():", c)
		panic(err)
	}
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// 加入链路追踪的配置
	if err := initTracer(bc.Trace.Endpoint); err != nil {

		panic(err)
	}

	// 创建 Kafka Writer
	kafkaWriter := favkafka.NewLogKafkaWriter(bc.Kafka)
	defer kafkaWriter.Writer.Close()

	// 创建 Logger
	logger := log.With(log.NewStdLogger(kafkaWriter),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)

	app, cleanup, err := wireApp(bc.Kafka, bc.Service, bc.Prometheus, bc.Etcd, bc.Auth, bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
