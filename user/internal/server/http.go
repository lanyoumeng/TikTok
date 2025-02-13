package server

import (
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	v1 "user/api/user/v1"
	"user/internal/conf"
	"user/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.UserService, Auth *conf.Auth, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),

			selector.Server( //jwt鉴权 白名单
				jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
					return []byte(Auth.JwtKey), nil
				}, jwt.WithSigningMethod(jwtv5.SigningMethodHS256)),
			).Match(NewWhiteListMatcher()).Build(),

			tracing.Server(), //jaejer链路追踪
			metrics.Server( //监控 prometheus
				metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
				metrics.WithRequests(prom.NewCounter(_metricRequests)),
			),
		),
	}

	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterUserServiceHTTPServer(srv, greeter)
	return srv
}
