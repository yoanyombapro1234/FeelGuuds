package grpc

import (
	"context"
	"io"
	"log"
	"net"
	"time"

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

// dialer creates an in memory full duplex connection
func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	s := NewMockServer()
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

// MockGRPCService creates and returns a mock grpc service connection
func MockGRPCService(ctx context.Context) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

// NewMockServer creates a new mock server instance
func NewMockServer() *Server {
	config := &Config{
		Port:            9897,
		ServiceName:     "AuthenticationHandlerService",
		RpcDeadline:     100,
		RpcRetries:      5,
		RpcRetryTimeout: 100,
		RpcRetryBackoff: 10,
	}

	const serviceName string = "test"

	// initiate tracing engine
	tracerEngine, closer := InitializeTracingEngine(serviceName)
	defer closer.Close()
	ctx := context.Background()

	// initiate metrics engine
	serviceMetrics := InitializeMetricsEngine(serviceName)

	// initiate logging client
	logger := InitializeLoggingEngine(ctx)

	// authn client
	authnClient, err := InitializeAuthnClient(logger)
	if err != nil {
		logger.For(ctx).FatalM(err, "unable to setup mock server")
	}

	srv := &Server{
		config:        config,
		tracerEngine:  tracerEngine,
		metricsEngine: serviceMetrics.Engine,
		metrics:       serviceMetrics.MicroServiceMetrics,
		logger:        logger,
		authnClient:   authnClient,
	}

	return srv
}

// InitializeAuthnClient creates a connection to the authn service
func InitializeAuthnClient(logger core_logging.ILog) (*core_auth_sdk.Client, error) {
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
	}, constants.TEST_ORIGIN)

	// TODO: make this a retryable operation
	retries := 1
	for retries < 4 {
		// perform a test request to the authentication service
		data, err := client.ServerStats()
		if err != nil {
			if retries != 4 {
				logger.ErrorM(err, "failed to connect to authentication service")
			} else {
				logger.FatalM(err, "failed to connect to authentication service")
			}
			retries += 1
		} else {
			retries = 4
			logger.InfoM("data", zap.Any("result", data))
		}

		time.Sleep(1 * time.Second)
	}

	return client, err
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
