package grpc

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	logger *zap.Logger
	config *Config
}

type Config struct {
	Port        int    `mapstructure:"grpc-port"`
	ServiceName string `mapstructure:"grpc-service-name"`
}

func NewServer(config *Config, logger *zap.Logger) (*Server, error) {
	srv := &Server{
		logger: logger,
		config: config,
	}

	return srv, nil
}

func (s *Server) ListenAndServe() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.config.Port))
	if err != nil {
		s.logger.Fatal("failed to listen", zap.Int("port", s.config.Port))
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
	grpc_health_v1.RegisterHealthServer(srv, server)
	server.SetServingStatus(s.config.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	if err := srv.Serve(listener); err != nil {
		s.logger.Fatal("failed to serve", zap.Error(err))
	}
}
