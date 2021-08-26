package grpc

import (
	"crypto/tls"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-auth-sdk"
	core_metrics "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-metrics"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-tracing/jaeger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/metrics"
)

type Server struct {
	config               *Config
	authnClient          *core_auth_sdk.Client
	logger               *zap.Logger
	metrics              *metrics.CoreMetrics
	metricsEngine        *core_metrics.CoreMetricsEngine
	tracerEngine         *core_tracing.TracingEngine
	enableTls            bool
	cert                 *tls.Certificate
	authnServiceMockStub core_auth_sdk.AuthService
}

var _ proto.AuthenticationHandlerServiceApiServer = (*Server)(nil)

type Config struct {
	Port                        int    `mapstructure:"GRPC_PORT"`
	ServiceName                 string `mapstructure:"GRPC_SERVICE_NAME"`
	RpcDeadline                 int    `mapstructure:"GRPC_RPC_DEADLINE_IN_MS"`
	RpcRetries                  int    `mapstructure:"GRPC_RPC_RETRIES"`
	RpcRetryTimeout             int    `mapstructure:"GRPC_RPC_RETRY_TIMEOUT_IN_MS"`
	RpcRetryBackoff             int    `mapstructure:"GRPC_RPC_RETRY_BACKOFF_IN_MS"`
	EnableTls                   bool   `mapstructure:"GRPC_ENABLE_TLS"`
	CertificatePath             string `mapstructure:"GRPC_CERT_PATH"`
	EnableDelayMiddleware       bool   `mapstructure:"ENABLE_RANDOM_DELAY"`
	EnableRandomErrorMiddleware bool   `mapstructure:"ENABLE_RANDOM_RANDOM_ERROR"`
	MinRandomDelay              int    `mapstructure:"RANDOM_DELAY_MIN_IN_MS"`
	MaxRandomDelay              int    `mapstructure:"RANDOM_DELAY_MAX_IN_MS"`
	DelayUnit                   string `mapstructure:"RANDOM_DELAY_UNIT"`
	Version                     string `mapstructure:"VERSION"`
	MetricAddr                  string `mapstructure:"METRIC_CONNECTION_ADDRESS"`
}

// NewGRPCServer defines a new instance of the grpc service
func NewGRPCServer(config *Config, client *core_auth_sdk.Client, logging *zap.Logger, serviceMetrics *metrics.CoreMetrics,
	metricsEngineConf *core_metrics.CoreMetricsEngine, tracer *core_tracing.TracingEngine) (*Server, error) {
	srv := &Server{
		logger:               logging,
		metrics:              serviceMetrics,
		authnClient:          client,
		metricsEngine:        metricsEngineConf,
		tracerEngine:         tracer,
		config:               config,
		authnServiceMockStub: nil,
	}

	return srv, nil
}

// ListenAndServe starts the grpc service
func (s *Server) ListenAndServe() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.config.Port))
	if err != nil {
		var msg = fmt.Sprintf("faled to listen on port %d", s.config.Port)
		s.logger.Fatal(err.Error(), zap.String("msg", msg))
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
		s.logger.Fatal(err.Error(), zap.String("msg", msg))
	}
}

/*
// ListenAndServe starts the grpc service
func (s *Server) ListenAndServe() {
	params := core_grpc.GrpcServerConfigurations{
		Logger:                   s.logger,
		GrpcServerConnectionAddr: fmt.Sprintf(":%v", s.config.Port),
		EnableTls:                s.enableTls,
		ServiceConfigs: &core_middleware.Configurations{
			StatsDConnectionAddr:        s.config.MetricAddr,
			Logger:                      s.logger,
			Client:                      s.authnClient,
			ServiceName:                 s.config.ServiceName,
			Origins:                     nil,
			EnableDelayMiddleware:       s.config.EnableDelayMiddleware,
			EnableRandomErrorMiddleware: s.config.EnableRandomErrorMiddleware,
			MinDelay:                    s.config.MinRandomDelay,
			MaxDelay:                    s.config.MaxRandomDelay,
			DelayUnit:                   s.config.DelayUnit,
			Version:                     s.config.Version,
		},
	}

		builder := grpc_inst.GrpcServerBuilder{}
		core_grpc.AddInterceptors(&builder, params.ServiceConfigs)
		builder.EnableReflection(true)
		builder.SetTlsCert(&tlscert.Cert)
		sv := builder.Build()
		sv.RegisterService(s.ServiceRegistration)
		err := sv.Start(params.GrpcServerConnectionAddr)
		if err != nil {
			s.logger.Fatal( err.Error())
		}
		sv.AwaitTermination(func() {
			s.logger.Info("Shutting down the server")
		})


	grpc_inst.ServerInitialization()
	grpcServer := core_grpc.NewGrpcService(&params)
	grpcServer.StartGrpcServer(s.ServiceRegistration)
}
*/

func (s *Server) ServiceRegistration(sv *grpc.Server) {
	proto.RegisterAuthenticationHandlerServiceApiServer(sv, s)
}
