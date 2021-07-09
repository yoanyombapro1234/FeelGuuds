package grpc

import (
	"fmt"

	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-tracing"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	grpc_client "github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/grpc-client"
	grpc_utils "github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/grpc-utils"
	tlscert "github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/tlsCert"
	"google.golang.org/grpc"
)

type Server struct {
	// inherit the behaviors/adhere to the interface the api server adheres to
	config                      *Config
	logger                      core_logging.ILog
	tracerEngine                *core_tracing.TracingEngine
	AuthenticationHandlerClient proto.AuthenticationHandlerServiceApiClient
}

var _ merchant_service_proto_v1.MerchantServiceServer = (*Server)(nil)

type Config struct {
	Port                                int    `mapstructure:"GRPC_PORT"`
	ServiceName                         string `mapstructure:"GRPC_SERVICE_NAME"`
	RpcDeadline                         int    `mapstructure:"grpc-rpc-deadline"`
	RpcRetries                          int    `mapstructure:"grpc-rpc-retries"`
	RpcRetryTimeout                     int    `mapstructure:"grpc-rpc-retry-timeout"`
	RpcRetryBackoff                     int    `mapstructure:"grpc-rpc-retry-backoff"`
	AuthenticationHandlerServiceAddress string `mapstructure:"grpc-authentication_handler_service-addr"`
	PaymentServiceAddress               string `mapstructure:"grpc-payment_service-addr"`
}

func NewServer(config *Config, logging core_logging.ILog, tracer *core_tracing.TracingEngine) (*Server, error) {
	srv := &Server{
		logger:       logging,
		config:       config,
		tracerEngine: tracer,
	}

	return srv, nil
}

func (s *Server) ListenAndServe(enableTls bool) {
	builder := s.InitializeServiceBuilder()
	sb := builder.Build()
	sb.RegisterService(s.serviceRegister)

	if enableTls {
		builder.SetTlsCert(&tlscert.Cert)
	}

	authHandlerSvcConn, err := grpc_client.ConnectToClient(s.config.AuthenticationHandlerServiceAddress)
	if err != nil {
		s.logger.Fatal(err, err.Error())
	}

	defer func(authHandlerSvcConn *grpc.ClientConn) {
		err := authHandlerSvcConn.Close()
		if err != nil {
			s.logger.Error(err, err.Error())
		}
	}(authHandlerSvcConn)

	s.AuthenticationHandlerClient = proto.NewAuthenticationHandlerServiceApiClient(authHandlerSvcConn)
	s.StartGrpcServer(err, sb)
}

func (s *Server) StartGrpcServer(err error, sb GrpcServer) {
	err = sb.Start(fmt.Sprintf(":%v", s.config.Port))
	if err != nil {
		s.logger.Fatal(err, err.Error())
	}

	sb.AwaitTermination(func() {
		s.logger.Info("Shutting down the server")
	})
}

func (s *Server) InitializeServiceBuilder() *GrpcServerBuilder {
	builder := NewGrpcServerBuilder(s.config.ServiceName, s.logger)
	s.addInterceptors(builder)
	builder.EnableReflection(true)
	return builder
}

func (s *Server) serviceRegister(sv *grpc.Server) {
	merchant_service_proto_v1.RegisterMerchantServiceServer(sv, s)
}

func (s *Server) addInterceptors(sb *GrpcServerBuilder) {
	sb.SetUnaryInterceptors(grpc_utils.GetDefaultUnaryServerInterceptors(s.tracerEngine.Tracer))
	sb.SetStreamInterceptors(grpc_utils.GetDefaultStreamServerInterceptors(s.tracerEngine.Tracer))
}
