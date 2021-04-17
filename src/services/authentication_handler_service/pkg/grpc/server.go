package grpc

import (
	"fmt"
	"net"

	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-auth-sdk"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	core_metrics "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-metrics"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	proto "github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/metrics"
)

type Server struct {
	// inherit the behaviors/adhere to the interface the api server adheres to
	proto.UnimplementedAuthenticationHandlerServiceApiServer
	config        *Config
	authnClient   *core_auth_sdk.Client
	logger        core_logging.ILog
	metrics       *metrics.CoreMetrics
	metricsEngine *core_metrics.CoreMetricsEngine
	tracerEngine  *core_tracing.TracingEngine
}

type Config struct {
	Port        int    `mapstructure:"grpc-port"`
	ServiceName string `mapstructure:"grpc-service-name"`
}

func NewServer(config *Config, client *core_auth_sdk.Client, logging core_logging.ILog, serviceMetrics *metrics.CoreMetrics,
	metricsEngineConf *core_metrics.CoreMetricsEngine, tracer *core_tracing.TracingEngine) (*Server, error) {
	srv := &Server{
		logger:        logging,
		metrics:       serviceMetrics,
		metricsEngine: metricsEngineConf,
		tracerEngine:  tracer,
		config:        config,
	}

	return srv, nil
}

func (s *Server) ListenAndServe() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.config.Port))
	if err != nil {
		var msg = fmt.Sprintf("faled to listen on port %d", s.config.Port)
		s.logger.FatalM(err, msg)
	}

	srv := grpc.NewServer()
	server := health.NewServer()
	reflection.Register(srv)

	// use the auto generate code to register server
	proto.RegisterAuthenticationHandlerServiceApiServer(srv, s)
	grpc_health_v1.RegisterHealthServer(srv, server)
	server.SetServingStatus(s.config.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	if err := srv.Serve(listener); err != nil {
		var msg = fmt.Sprintf("faled to serve on port %d", s.config.Port)
		s.logger.FatalM(err, msg)
	}
}
