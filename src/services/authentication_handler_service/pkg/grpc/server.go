package grpc

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-auth-sdk"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	core_metrics "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-metrics"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-tracing"

	otgrpc "github.com/opentracing-contrib/go-grpc"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/metrics"
)

type Server struct {
	// inherit the behaviors/adhere to the interface the api server adheres to
	proto.UnimplementedAuthenticationHandlerServiceApiServer
	config        *Config
	authnClient   core_auth_sdk.AuthService
	logger        core_logging.ILog
	metrics       *metrics.CoreMetrics
	metricsEngine *core_metrics.CoreMetricsEngine
	tracerEngine  *core_tracing.TracingEngine
}

type Config struct {
	Port            int    `mapstructure:"grpc-port"`
	ServiceName     string `mapstructure:"grpc-service-name"`
	RpcDeadline     int    `mapstructure:"grpc-rpc-deadline"`
	RpcRetries      int    `mapstructure:"grpc-rpc-retries"`
	RpcRetryTimeout int    `mapstructure:"grpc-rpc-retry-timeout"`
	RpcRetryBackoff int    `mapstructure:"grpc-rpc-retry-backoff"`
}

// NewServer defines a new instance of the grpc service
func NewGRPCServer(config *Config, client core_auth_sdk.AuthService, logging core_logging.ILog, serviceMetrics *metrics.CoreMetrics,
	metricsEngineConf *core_metrics.CoreMetricsEngine, tracer *core_tracing.TracingEngine) (*Server, error) {
	srv := &Server{
		logger:        logging,
		metrics:       serviceMetrics,
		authnClient:   client,
		metricsEngine: metricsEngineConf,
		tracerEngine:  tracer,
		config:        config,
	}

	return srv, nil
}

// ListenAndServe starts the grpc service
func (s *Server) ListenAndServe() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.config.Port))
	if err != nil {
		var msg = fmt.Sprintf("faled to listen on port %d", s.config.Port)
		s.logger.Fatal(err, msg)
	}

	// configure tracing so all future rpc activity will be traced by use of s
	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
			otgrpc.OpenTracingStreamServerInterceptor(s.tracerEngine.Tracer),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
			otgrpc.OpenTracingServerInterceptor(s.tracerEngine.Tracer)),
		))

	server := health.NewServer()
	reflection.Register(srv)

	// use the auto generate code to register server
	proto.RegisterAuthenticationHandlerServiceApiServer(srv, s)
	grpc_health_v1.RegisterHealthServer(srv, server)
	server.SetServingStatus(s.config.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	if err := srv.Serve(listener); err != nil {
		var msg = fmt.Sprintf("faled to serve on port %d", s.config.Port)
		s.logger.Fatal(err, msg)
	}
}
