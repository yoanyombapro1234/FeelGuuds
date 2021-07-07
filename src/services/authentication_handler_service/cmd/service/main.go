package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"

	core_metrics "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-metrics"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-tracing"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/api"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/grpc"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/metrics"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/signals"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/version"

	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-auth-sdk"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
)

func main() {
	// flags definition
	fs := pflag.NewFlagSet("default", pflag.ContinueOnError)
	fs.Int("port", 9898, "HTTP port")
	fs.Int("secure-port", 0, "HTTPS port")
	fs.Int("port-metrics", 0, "metrics port")
	fs.Int("grpc-port", 9897, "gRPC port")
	fs.String("grpc-service-name", "service", "gPRC service name")
	fs.Int("grpc-rpc-deadline", 500, "gRPC deadline in milliseconds")
	fs.Int("grpc-rpc-retries", 1, "gRPC max operation retries in the face of errors")
	fs.Int("grpc-rpc-retry-timeout", 100, "gRPC max timeout of retry operation in milliseconds")
	fs.Int("grpc-rpc-retry-backoff", 20, "gRPC backoff in between failed retry operations in milliseconds")

	fs.String("level", "info", "log level debug, info, warn, error, flat or panic")
	fs.StringSlice("backend-url", []string{}, "backend service URL")
	fs.Duration("http-client-timeout", 2*time.Minute, "client timeout duration")
	fs.Duration("http-server-timeout", 30*time.Second, "server read and write timeout duration")
	fs.Duration("http-server-shutdown-timeout", 5*time.Second, "server graceful shutdown timeout duration")
	fs.String("data-path", "/data", "data local path")
	fs.String("config-path", "", "config dir path")
	fs.String("cert-path", "/data/cert", "certificate path for HTTPS port")
	fs.String("config", "config.yaml", "config file name")
	fs.String("ui-path", "./ui", "UI local path")
	fs.String("ui-logo", "", "UI logo")
	fs.String("ui-color", "#34577c", "UI color")
	fs.String("ui-message", fmt.Sprintf("greetings from service v%v", version.VERSION), "UI message")
	fs.Bool("h2c", false, "allow upgrading to H2C")
	fs.Bool("random-delay", false, "between 0 and 5 seconds random delay by default")
	fs.String("random-delay-unit", "s", "either s(seconds) or ms(milliseconds")
	fs.Int("random-delay-min", 0, "min for random delay: 0 by default")
	fs.Int("random-delay-max", 5, "max for random delay: 5 by default")
	fs.Bool("random-error", false, "1/3 chances of a random response error")
	fs.Bool("unhealthy", false, "when set, healthy state is never reached")
	fs.Bool("unready", false, "when set, ready state is never reached")
	fs.Int("stress-cpu", 0, "number of CPU cores with 100 load")
	fs.Int("stress-memory", 0, "MB of data to load into memory")
	fs.String("cache-server", "", "Redis address in the format <host>:<port>")

	// authentication service specific flags
	fs.String("SERVICE_AUTHN_USERNAME", "feelguuds", "username of authentication client")
	fs.String("SERVICE_AUTHN_PASSWORD", "feelguuds", "password of authentication client")
	fs.String("SERVICE_AUTHN_ISSUER_BASE_URL", "http://localhost", "authentication service issuer")
	fs.String("SERVICE_AUTHN_ORIGIN", "http://localhost", "origin of auth requests")
	fs.String("SERVICE_AUTHN_DOMAINS", "localhost", "authentication service domains")
	fs.String("SERVICE_AUTHN_PRIVATE_BASE_URL", "http://authentication_service",
		"authentication service private url. should be local host if these are not running on docker containers. "+
			"However if running in docker container with a configured docker network, the url should be equal to the service name")
	fs.String("SERVICE_AUTHN_PUBLIC_BASE_URL", "http://localhost", "authentication service public endpoint")
	fs.String("SERVICE_AUTHN_INTERNAL_PORT", "3000", "authentication service port")
	fs.String("SERVICE_AUTHN_PORT", "8000", "authentication service external port")
	fs.Bool("SERVICE_ENABLE_AUTH_SERVICE_PRIVATE_INTEGRATION", true, "enables communication with authentication service")

	// retry specific configurations
	fs.Int("HTTP_MAX_RETRIES", 3, "max retries to perform on failed http calls")
	fs.Duration("HTTP_MIN_RETRY_WAITING_TIME", 10*time.Millisecond, "minimum time to wait between failed calls for retry")
	fs.Duration("HTTP_MAX_RETRY_WAITING_TIME", 50*time.Millisecond, "maximum time to wait between failed calls for retry")
	fs.Duration("HTTP_REQUEST_TIMEOUT", 5*time.Second, "time until a request is seen as timing out")

	// logging specific configurations
	fs.String("SERVICE_NAME", "authentication_handler_service", "service name")
	// TODO: reconfigure this to leverage datadog instead
	fs.String("JAEGER_ENDPOINT", "http://jaeger-collector:14268/api/traces", "jaeger collector endpoint")

	// capture goroutines waiting on synchronization primitives
	runtime.SetBlockProfileRate(1)

	versionFlag := fs.BoolP("version", "v", false, "get version number")

	// parse flags
	err := fs.Parse(os.Args[1:])
	switch {
	case err == pflag.ErrHelp:
		os.Exit(0)
	case err != nil:
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		fs.PrintDefaults()
		os.Exit(2)
	case *versionFlag:
		fmt.Println(version.VERSION)
		os.Exit(0)
	}

	// bind flags and environment variables
	viper.BindPFlags(fs)
	viper.RegisterAlias("backendUrl", "backend-url")
	hostname, _ := os.Hostname()
	viper.SetDefault("jwt-secret", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	viper.SetDefault("ui-logo", "https://raw.githubusercontent.com/stefanprodan/podinfo/gh-pages/cuddle_clap.gif")
	viper.Set("hostname", hostname)
	viper.Set("version", version.VERSION)
	viper.Set("revision", version.REVISION)
	viper.SetEnvPrefix("SERVICE")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// load config from file
	if _, fileErr := os.Stat(filepath.Join(viper.GetString("config-path"), viper.GetString("config"))); fileErr == nil {
		viper.SetConfigName(strings.Split(viper.GetString("config"), ".")[0])
		viper.AddConfigPath(viper.GetString("config-path"))
		if readErr := viper.ReadInConfig(); readErr != nil {
			fmt.Printf("Error reading config file, %v\n", readErr)
		}
	} else {
		fmt.Printf("Error to open config file, %v\n", fileErr)
	}

	serviceName := viper.GetString("SERVICE_NAME")
	collectorEndpoint := viper.GetString("JAEGER_ENDPOINT")

	// initialize a tracing object globally
	tracerEngine, closer := core_tracing.NewTracer(serviceName, collectorEndpoint, prometheus.New())
	defer closer.Close()

	if tracerEngine == nil {
		panic("cannot initialize tracer engine")
	}
	opentracing.SetGlobalTracer(tracerEngine.Tracer)

	// initialize metrics object
	coreMetrics := core_metrics.NewCoreMetricsEngineInstance(serviceName, nil)
	serviceMetrics := metrics.NewMetricsEngine(coreMetrics, serviceName)

	// start root span
	ctx := context.Background()
	rootSpan := opentracing.SpanFromContext(ctx)

	// configure logging
	logger := core_logging.NewJSONLogger(nil, rootSpan)

	authnServiceClient := NewAuthServiceClientConnection(err, logger)
	if authnServiceClient != nil {
		logger.Info("successfully initialized authentication service client")
	}

	// start stress tests if any
	beginStressTest(viper.GetInt("stress-cpu"), viper.GetInt("stress-memory"), logger)

	// validate port
	if _, err := strconv.Atoi(viper.GetString("port")); err != nil {
		port, _ := fs.GetInt("port")
		viper.Set("port", strconv.Itoa(port))
	}

	// validate secure port
	if _, err := strconv.Atoi(viper.GetString("secure-port")); err != nil {
		securePort, _ := fs.GetInt("secure-port")
		viper.Set("secure-port", strconv.Itoa(securePort))
	}

	// validate random delay options
	if viper.GetInt("random-delay-max") < viper.GetInt("random-delay-min") {
		err := errors.New("`--random-delay-max` should be greater than `--random-delay-min`")
		logger.Fatal(err, "please fix configurations")
	}

	switch delayUnit := viper.GetString("random-delay-unit"); delayUnit {
	case
		"s",
		"ms":
		break
	default:
		err := errors.New("random-delay-unit` accepted values are: s|ms")
		logger.Fatal(err, "please fix configurations")
	}

	// load gRPC server config
	var grpcCfg grpc.Config
	if err := viper.Unmarshal(&grpcCfg); err != nil {
		err := errors.New("config unmarshal failed")
		logger.Fatal(err, "please fix configurations")
	}

	// start gRPC server
	if grpcCfg.Port > 0 {
		logger.Info("starting grpc server")
		grpcSrv, _ := grpc.NewGRPCServer(&grpcCfg, authnServiceClient, logger, serviceMetrics.MicroServiceMetrics, serviceMetrics.Engine, tracerEngine)

		logger.Info("successfully started grpc server", zap.Int("port", grpcCfg.Port))
		go grpcSrv.ListenAndServe()
	}

	// load HTTP server config
	var srvCfg api.Config
	if err := viper.Unmarshal(&srvCfg); err != nil {
		logger.Fatal(err, "config unmarshal failed")
	}

	// log version and port
	logger.Info("Starting service",
		zap.String("version", viper.GetString("version")),
		zap.String("revision", viper.GetString("revision")),
		zap.String("port", srvCfg.Port),
	)

	// start HTTP server
	srv, _ := api.NewServer(&srvCfg, authnServiceClient, logger, serviceMetrics.MicroServiceMetrics, serviceMetrics.Engine, tracerEngine)
	stopCh := signals.SetupSignalHandler()
	srv.ListenAndServe(stopCh)
}

var stressMemoryPayload []byte

func beginStressTest(cpus int, mem int, logger core_logging.ILog) {
	done := make(chan int)
	if cpus > 0 {
		logger.Info("starting CPU stress", zap.Int("cores", cpus))
		for i := 0; i < cpus; i++ {
			go func() {
				for {
					select {
					case <-done:
						return
					default:

					}
				}
			}()
		}
	}

	if mem > 0 {
		path := "/tmp/service.data"
		f, err := os.Create(path)

		if err != nil {
			logger.Error(err, "memory stress failed", "error")
		}

		if err := f.Truncate(1000000 * int64(mem)); err != nil {
			logger.Error(err, "memory stress failed", "error")
		}

		stressMemoryPayload, err = ioutil.ReadFile(path)
		f.Close()
		os.Remove(path)
		if err != nil {
			logger.Error(err, "memory stress failed", "error")
		}
		logger.Info("starting CPU stress", zap.Int("memory", len(stressMemoryPayload)))
	}
}

// initAuthnClient initializes an instance of the authn client primarily useful in
// communicating with the authentication service securely
func initAuthnClient(username, password, audience, issuer, url, origin string) (*core_auth_sdk.Client, error) {
	retryConfig := &core_auth_sdk.RetryConfig{
		MaxRetries:       viper.GetInt("HTTP_MAX_RETRIES"),
		MinRetryWaitTime: viper.GetDuration("HTTP_MIN_RETRY_WAITING_TIME"),
		MaxRetryWaitTime: viper.GetDuration("HTTP_MAX_RETRY_WAITING_TIME"),
		RequestTimeout:   viper.GetDuration("HTTP_REQUEST_TIMEOUT"),
	}

	// Authentication.
	return core_auth_sdk.NewClient(core_auth_sdk.Config{
		// The AUTHN_URL of your Keratin AuthN server. This will be used to verify tokens created by
		// AuthN, and will also be used for API calls unless PrivateBaseURL is also set.
		Issuer: issuer,

		// The domain of your application (no protocol). This domain should be listed in the APP_DOMAINS
		// of your Keratin AuthN server.
		Audience: audience,

		// Credentials for AuthN's private endpoints. These will be used to execute admin actions using
		// the Client provided by this library.
		//
		// TIP: make them extra secure in production!
		Username: username,
		Password: password,

		// RECOMMENDED: Send private API calls to AuthN using private network routing. This can be
		// necessary if your environment has a firewall to limit public endpoints.
		PrivateBaseURL: url,
	}, origin, retryConfig)
}

func NewAuthServiceClientConnection(err error, logger core_logging.ILog) core_auth_sdk.AuthService {
	// initialize authentication client in order to establish communication with the
	// authentication service. This serves as a singular source of truth for authentication needs
	authUsername := viper.GetString("SERVICE_AUTHN_USERNAME")
	authPassword := viper.GetString("SERVICE_AUTHN_PASSWORD")
	domains := viper.GetString("SERVICE_AUTHN_DOMAINS")
	privateURL := viper.GetString("SERVICE_AUTHN_PRIVATE_BASE_URL") + ":" + viper.GetString("SERVICE_AUTHN_INTERNAL_PORT")
	origin := viper.GetString("SERVICE_AUTHN_ORIGIN")
	issuer := viper.GetString("SERVICE_AUTHN_ISSUER_BASE_URL") + ":" + viper.GetString("SERVICE_AUTHN_PORT")

	authnClient, err := initAuthnClient(authUsername, authPassword, domains, issuer, privateURL, origin)
	// crash the process if we cannot connect to the authentication service
	if err != nil {
		logger.Fatal(err, "failed to initialized authentication service client")
	}

	// TODO: make this a retryable operation
	retries := 1
	retryLimit := 8
	for retries < retryLimit {
		// perform a test request to the authentication service
		_, err = authnClient.ServerStats()
		if err != nil {
			if retries != retryLimit {
				logger.Error(err, fmt.Sprintf("failed to connect to authentication service. Attempt #%d", retries))
			}
			retries += 1
		} else {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if err != nil {
		logger.Error(errors.New("failed to initiate connection to downstream service"), "failure")
		return nil
	}

	// attempt to connect to the authentication service if not then crash process
	return authnClient
}
