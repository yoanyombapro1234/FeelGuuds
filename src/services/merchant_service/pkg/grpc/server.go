package grpc

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-tracing"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	// inherit the behaviors/adhere to the interface the api server adheres to
	config        *Config
	logger        core_logging.ILog
	tracerEngine  *core_tracing.TracingEngine
}

var _ merchant_service_proto_v1.MerchantServiceServer = (*Server)(nil)

type Config struct {
	Port        int    `mapstructure:"GRPC_PORT"`
	ServiceName string `mapstructure:"GRPC_SERVICE_NAME"`
	RpcDeadline     int    `mapstructure:"grpc-rpc-deadline"`
	RpcRetries      int    `mapstructure:"grpc-rpc-retries"`
	RpcRetryTimeout int    `mapstructure:"grpc-rpc-retry-timeout"`
	RpcRetryBackoff int    `mapstructure:"grpc-rpc-retry-backoff"`
}

func NewServer(config *Config, logging core_logging.ILog, tracer *core_tracing.TracingEngine) (*Server, error) {
	srv := &Server{
		logger: logging,
		config: config,
		tracerEngine: tracer,
	}

	return srv, nil
}

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

	// use auto generate code to register server
	merchant_service_proto_v1.RegisterMerchantServiceServer(srv, s)
	grpc_health_v1.RegisterHealthServer(srv, server)
	server.SetServingStatus(s.config.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	if err := srv.Serve(listener); err != nil {
		var msg = fmt.Sprintf("faled to serve on port %d", s.config.Port)
		s.logger.Fatal(err, msg)
	}
}
