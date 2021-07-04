package grpc

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"github.com/giantswarm/retry-go"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-auth-sdk"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	core_metrics "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-metrics"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-tracing"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/metrics"
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
func NewMockServer(authClientMockStub core_auth_sdk.AuthService) *Server {
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
	ctx := context.Background()

	// initiate metrics engine
	serviceMetrics := InitializeMetricsEngine(config.ServiceName)

	// initiate logging client
	logger := InitializeLoggingEngine(ctx)

	if authClientMockStub == nil {
		authClientMockStub, err = InitializeAuthnClient(logger)
		if err != nil {
			logger.Fatal(err, err.Error())
		}
	}

	srv := &Server{
		config:        config,
		tracerEngine:  tracerEngine,
		metricsEngine: serviceMetrics.Engine,
		metrics:       serviceMetrics.MicroServiceMetrics,
		logger:        logger,
		authnClient:   authClientMockStub,
	}

	return srv
}

// InitializeAuthnClient creates a connection to the authn service
func InitializeAuthnClient(logger core_logging.ILog) (core_auth_sdk.AuthService, error) {
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
	retries := 1
	for retries < 4 {
		// perform a test request to the authentication service
		data, err := client.ServerStats()
		if err != nil {
			if retries != 4 {
				logger.Error(err, "failed to connect to authentication service")
			} else {
				logger.Fatal(err, "failed to connect to authentication service")
			}
			retries += 1
		} else {
			retries = 4
			logger.Info("data", zap.Any("result", data))
		}

		time.Sleep(3 * time.Second)
	}

	// err = ConnectToDownstreamService(logger, client, response)

	return client, err
}

// ConnectToDownstreamService attempts to connect to a downstream service
func ConnectToDownstreamService(logger core_logging.ILog, client *core_auth_sdk.Client, response chan interface{}) error {
	return retry.Do(
		func(conn chan<- interface{}) func() error {
			return func() error {
				data, err := client.ServerStats()
				if err != nil {
					logger.Error(err, "failed to connect to authentication service")
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
func InitializeLoggingEngine(ctx context.Context) core_logging.ILog {
	rootSpan := opentracing.SpanFromContext(ctx)
	logger := core_logging.NewJSONLogger(nil, rootSpan)
	return logger
}

// InitializeMetricsEngine initializes a metrics engine globally
func InitializeMetricsEngine(serviceName string) *metrics.MetricsEngine {
	coreMetrics := core_metrics.NewCoreMetricsEngineInstance(serviceName, nil)
	serviceMetrics := metrics.NewMetricsEngine(coreMetrics, "mock")
	return serviceMetrics
}

// InitializeTracingEngine initiaize a tracing object globally
func InitializeTracingEngine(serviceName string) (*core_tracing.TracingEngine, io.Closer) {
	return core_tracing.NewTracer(serviceName, constants.COLLECTOR_ENDPOINT, prometheus.New())
}
