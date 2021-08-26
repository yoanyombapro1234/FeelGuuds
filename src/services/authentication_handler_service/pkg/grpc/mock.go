package grpc

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"github.com/giantswarm/retry-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/metrics"
	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-auth-sdk"
	core_logging "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-logging"
	core_metrics "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-metrics"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-tracing/jaeger"
)

type MockAuthGRPCServer struct {
	proto.UnimplementedAuthenticationHandlerServiceApiServer
}

type MockDialOption func(context.Context, string) (net.Conn, error)

// dialer creates an in memory full duplex connection
func dialer(authClientMock core_auth_sdk.AuthService) func() MockDialOption {
	return func() MockDialOption {
		listener := bufconn.Listen(1024 * 1024)

		server := grpc.NewServer()
		s := NewMockServer(authClientMock)
		proto.RegisterAuthenticationHandlerServiceApiServer(server, s)

		go func() {
			if err := server.Serve(listener); err != nil {
				log.Fatal(err)
			}
		}()

		return func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}
	}
}

// MockGRPCService creates and returns a mock grpc service connection
func MockGRPCService(ctx context.Context, authClientMock core_auth_sdk.AuthService) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(authClientMock)()))
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

// NewMockServer creates a new mock server instance
func NewMockServer(authServiceMockStub core_auth_sdk.AuthService) *Server {
	var err error

	config := &Config{
		Port:            9897,
		ServiceName:     "AuthenticationHandlerService",
		RpcDeadline:     60 * 1000000,
		RpcRetries:      1,
		RpcRetryTimeout: 10,
		RpcRetryBackoff: 1,
	}

	// initiate tracing engine
	tracerEngine, closer := InitializeTracingEngine(config.ServiceName)
	defer closer.Close()

	// initiate metrics engine
	metricsEngine, serviceMetrics := InitializeMetricsEngine(config.ServiceName)

	// initiate logging client
	logger := InitializeLoggingEngine()

	if authServiceMockStub == nil {
		authServiceMockStub, err = InitializeAuthnClient(logger)
		if err != nil {
			logger.Fatal(err.Error())
		}
	}

	srv := &Server{
		config:               config,
		tracerEngine:         tracerEngine,
		metricsEngine:        metricsEngine,
		metrics:              serviceMetrics,
		logger:               logger,
		authnClient:          nil,
		authnServiceMockStub: authServiceMockStub,
	}

	return srv
}

// InitializeAuthnClient creates a connection to the authn service
func InitializeAuthnClient(logger *zap.Logger) (core_auth_sdk.AuthService, error) {
	// var response = make(chan interface{}, 1)

	client, err := core_auth_sdk.NewClient(core_auth_sdk.Config{
		// The AUTHN_URL of your Keratin AuthN server. This will be used to verify tokens created by
		// AuthN, and will also be used for API calls unless PrivateBaseURL is also set.
		Issuer: constants.TEST_ISSUER,

		// The domain of your application (no protocol). This domain should be listed in the APP_DOMAINS
		// of your Keratin AuthN server.
		Audience: constants.TEST_AUDIENCE,

		// Credentials for AuthN's private endpoints. These will be used to execute admin actions using
		// the Client provided by this library.
		//
		// TIP: make them extra secure in production!
		Username: constants.TEST_USERNAME,
		Password: constants.TEST_PASSWORD,

		// RECOMMENDED: Send private API calls to AuthN using private network routing. This can be
		// necessary if your environment has a firewall to limit public endpoints.
		PrivateBaseURL: constants.TEST_BASE_URL,
	}, constants.TEST_ORIGIN, &core_auth_sdk.RetryConfig{
		MaxRetries:       5,
		MinRetryWaitTime: 5 * time.Millisecond,
		MaxRetryWaitTime: 10 * time.Millisecond,
		RequestTimeout:   400 * time.Millisecond,
	})

	var response = make(chan interface{}, 1)
	err = ConnectToDownstreamService(logger, client, response)
	return client, err
}

// ConnectToDownstreamService attempts to connect to a downstream service
func ConnectToDownstreamService(logger *zap.Logger, client *core_auth_sdk.Client, response chan interface{}) error {
	return retry.Do(
		func(conn chan<- interface{}) func() error {
			return func() error {
				data, err := client.ServerStats()
				if err != nil {
					logger.Error("failed to connect to authentication service", zap.Error(err))
					return err
				}

				logger.Info("data", zap.Any("result", data))

				response <- data
				return nil
			}
		}(response),
		retry.MaxTries(5),
		retry.Timeout(time.Millisecond*time.Duration(10)),
		retry.Sleep(time.Millisecond*time.Duration(10)))
}

// InitializeLoggingEngine initializes logging object
func InitializeLoggingEngine() *zap.Logger {
	return core_logging.New("info").Logger
}

// InitializeMetricsEngine initializes a metrics engine globally
func InitializeMetricsEngine(serviceName string) (*core_metrics.CoreMetricsEngine, *metrics.CoreMetrics) {
	coreMetrics := core_metrics.NewCoreMetricsEngineInstance(serviceName, nil)
	serviceMetrics := metrics.New(coreMetrics, "mock")
	return coreMetrics, serviceMetrics.MicroServiceMetrics
}

// InitializeTracingEngine initialize a tracing object globally
func InitializeTracingEngine(serviceName string) (*core_tracing.TracingEngine, io.Closer) {
	return core_tracing.New(serviceName, constants.COLLECTOR_ENDPOINT)
}
