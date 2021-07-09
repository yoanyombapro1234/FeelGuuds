package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-tracing"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/database"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	grpc_client "github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/grpc-client"
	grpc_utils "github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/grpc-utils"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/stripe_client"
	tlscert "github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/tlsCert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	// inherit the behaviors/adhere to the interface the api server adheres to
	config                      *Config
	logger                      core_logging.ILog
	tracerEngine                *core_tracing.TracingEngine
	AuthenticationHandlerClient proto.AuthenticationHandlerServiceApiClient
	DbConn                      *database.Db
	StripeClient                *stripe_client.Client
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
	StripeKey                           string `mapstructure:"stripe-key"`
	RefreshUrl                          string `mapstructure:"refresh-url"`
	ReturnUrl                           string `mapstructure:"return-url"`
}

func NewServer(config *Config, logging core_logging.ILog, tracer *core_tracing.TracingEngine, dbConn *database.Db) (*Server, error) {
	if config == nil || dbConn == nil || tracer == nil {
		return nil, errors.ErrInvalidInputArguments
	}

	client := stripe_client.NewStripeClient(logging, stripe_client.ClientParams{
		Key:        config.StripeKey,
		RefreshUrl: config.RefreshUrl,
		ReturnUrl:  config.ReturnUrl,
	})

	srv := &Server{
		logger:       logging,
		config:       config,
		tracerEngine: tracer,
		DbConn:       dbConn,
		StripeClient: client,
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

// ConfigureAndStartRootSpan configures a parent span object and starts it
func (s *Server) ConfigureAndStartRootSpan(ctx context.Context, operationType string) (context.Context, opentracing.Span) {
	ctx, _ = s.setCtxRequestTimeout(ctx)
	ctx, rootSpan := s.StartRootSpan(ctx, operationType)
	return ctx, rootSpan
}

// setCtxRequestTimeout sets the request deadline in the context. This function should be invoked prior to any rpc calls
func (s *Server) setCtxRequestTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	clientDeadline := time.Now().Add(time.Duration(s.config.RpcDeadline) * time.Millisecond)
	return context.WithDeadline(ctx, clientDeadline)
}

// StartRootSpan starts the rootspan of the current operation at hand
func (s *Server) StartRootSpan(ctx context.Context, operationType string) (context.Context, opentracing.Span) {
	s.logger.For(ctx).Info("GRPC request received", zap.String("method", operationType))

	spanCtx, _ := s.tracerEngine.Tracer.Extract(opentracing.HTTPHeaders, nil)
	parentSpan := s.tracerEngine.Tracer.StartSpan(operationType, ext.RPCServerOption(spanCtx))
	ctx = opentracing.ContextWithSpan(ctx, parentSpan)
	return ctx, parentSpan
}
