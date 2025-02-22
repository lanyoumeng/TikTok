package server

import (
	"context"
	v1 "message/api/message/v1"
	"message/internal/conf"
	"message/internal/service"

	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	_metricSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "duration_sec",
		Help:      "server requests duratio(sec).",
		Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1},
	}, []string{"kind", "operation"})

	_metricRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "client",
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "The total number of processed requests",
	}, []string{"kind", "operation", "code", "reason"})
)

func init() {
	prometheus.MustRegister(_metricSeconds, _metricRequests)
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(pro *conf.Prometheus, c *conf.Server, Auth *conf.Auth, message *service.MessageService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),

			//selector.Server( //中间件 权限验证
			//	jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
			//		return []byte(Auth.JwtKey), nil
			//	}, jwt.WithSigningMethod(jwtv5.SigningMethodHS256)),
			//).Match(NewWhiteListMatcher()).Build(),

			tracing.Server(), //jaejer链路追踪
			metrics.Server( //监控 prometheus
				metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
				metrics.WithRequests(prom.NewCounter(_metricRequests)),
			),
		),
	}

	go func() {
		//暴露 prometheus 监控指标
		addr := pro.Addr
		path := pro.Path
		httpSrv := http.NewServer(
			http.Address(addr),
			http.Middleware(
				metrics.Server(
					metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
					metrics.WithRequests(prom.NewCounter(_metricRequests)),
				),
			),
		)
		httpSrv.Handle(path, promhttp.Handler())
		if err := httpSrv.Start(context.Background()); err != nil {
			panic(err)
		}

	}()

	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterMessageServiceServer(srv, message)
	//log.NewHelper(logger).Debug(" grpc.NewServer", srv)

	return srv
}

func NewWhiteListMatcher() selector.MatchFunc {

	whiteList := make(map[string]struct{})
	whiteList["/message.api.message.v1.MessageService/MessageRecord"] = struct{}{}
	whiteList["/message.api.message.v1.MessageService/MessageSend"] = struct{}{}
	whiteList["/message.api.message.v1.MessageService/GetNewMessages"] = struct{}{}

	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}
