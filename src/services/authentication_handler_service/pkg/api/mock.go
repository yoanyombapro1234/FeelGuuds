package api

import (
	"io"
	"time"

	"github.com/giantswarm/retry-go"
	"github.com/gorilla/mux"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/metrics"
	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-auth-sdk"
	core_logging "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-logging"
	core_metrics "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-metrics"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-tracing/jaeger"
)

const collectorEndpoint string = "http://localhost:14268/api/traces"

func NewMockServer() *Server {
	config := &Config{
		Port:                      "9898",
		HttpServerShutdownTimeout: 5 * time.Second,
		HttpServerTimeout:         30 * time.Second,
		BackendURL:                []string{},
		ConfigPath:                "/config",
		DataPath:                  "/data",
		HttpClientTimeout:         30 * time.Second,
		UIColor:                   "blue",
		UIPath:                    ".ui",
		UIMessage:                 "Greetings",
		Hostname:                  "localhost",
	}

	const serviceName string = "test"

	// initiate tracing engine
	tracerEngine, closer := InitializeTracingEngine(serviceName)
	defer closer.Close()

	// initiate metrics engine
	metricsEngine, serviceMetrics := InitializeMetricsEngine(serviceName)

	// initiate logging client
	logger := InitializeLoggingEngine()

	// authn client
	authnClient, err := InitializeAuthnClient(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	srv := &Server{
		router:        mux.NewRouter(),
		config:        config,
		tracerEngine:  tracerEngine,
		metricsEngine: metricsEngine,
		metrics:       serviceMetrics,
		logger:        logger,
		authnClient:   authnClient,
	}
	return srv
}

func InitializeAuthnClient(logger *zap.Logger) (*core_auth_sdk.Client, error) {
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
		MaxRetryWaitTime: 15 * time.Millisecond,
		RequestTimeout:   400 * time.Millisecond,
	})

	var response = make(chan interface{}, 1)
	err = retry.Do(
		func(conn chan<- interface{}) func() error {
			return func() error {
				opResponse, err := client.ServerStats()
				if err != nil {
					return err
				}
				response <- opResponse
				return nil
			}
		}(response),
		retry.MaxTries(5),
		retry.Timeout(time.Millisecond*time.Duration(500)),
		retry.Sleep(time.Millisecond*time.Duration(50)),
	)

	if err != nil {
		logger.Error(err.Error())
	}
	return client, err
}

func InitializeLoggingEngine() *zap.Logger {
	logger := core_logging.New("info")
	return logger.Logger
}

func InitializeMetricsEngine(serviceName string) (*core_metrics.CoreMetricsEngine, *metrics.CoreMetrics) {
	coreMetrics := core_metrics.NewCoreMetricsEngineInstance(serviceName, nil)
	serviceMetrics := metrics.New(coreMetrics, "mock")
	return coreMetrics, serviceMetrics.MicroServiceMetrics
}

func InitializeTracingEngine(serviceName string) (*core_tracing.TracingEngine, io.Closer) {
	return core_tracing.New(serviceName, collectorEndpoint)
}
