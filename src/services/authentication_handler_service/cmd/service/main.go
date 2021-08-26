package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/giantswarm/retry-go"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/api"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/grpc"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/metrics"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/signals"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/version"
	core_metrics "github.com/yoanyombapro1234/FeelGuuds_core/core/core-metrics"
	tracer "github.com/yoanyombapro1234/FeelGuuds_core/core/core-tracing/jaeger"

	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_core/core/core-auth-sdk"
	core_logging "github.com/yoanyombapro1234/FeelGuuds_core/core/core-logging"
)

func main() {
	// flags definition
	fs := pflag.NewFlagSet("default", pflag.ContinueOnError)
	fs.Int("HTTP_PORT", 9897, "HTTP port")
	fs.Int("HTTPS_PORT", 9898, "HTTPS port")
	fs.Int("METRICS_PORT", 9899, "metrics port")
	fs.Int("GRPC_PORT", 9896, "gRPC port")
	fs.String("GRPC_SERVICE_NAME", "AUTHENTICATION_HANDLER_SERVICE", "gPRC service name")
	fs.Int("GRPC_RPC_DEADLINE_IN_MS", 5, "gRPC deadline in milliseconds")
	fs.Int("GRPC_RPC_RETRIES", 2, "gRPC max operation retries in the face of errors")
	fs.Int("GRPC_RPC_RETRY_TIMEOUT_IN_MS", 100, "gRPC max timeout of retry operation in milliseconds")
	fs.Int("GRPC_RPC_RETRY_BACKOFF_IN_MS", 20, "gRPC backoff in between failed retry operations in milliseconds")

	fs.String("LOG_LEVEL", "info", "log level debug, info, warn, error, flat or panic")
	fs.StringSlice("BACKEND_SERVICE_URLS", []string{}, "backend service URL")
	fs.Duration("HTTP_CLIENT_TIMEOUT_IN_MINUTES", 2*time.Minute, "client timeout duration")
	fs.Duration("HTTP_SERVER_TIMEOUT_IN_SECONDS", 30*time.Second, "server read and write timeout duration")
	fs.Duration("HTTP_SERVER_SHUTDOWN_TIMEOUT_IN_SECONDS", 5*time.Second, "server graceful shutdown timeout duration")
	fs.String("DATA_PATH", "/data", "data local path")
	fs.String("CONFIG_PATH", "", "config dir path")
	fs.String("CERT_PATH", "/data/cert", "certificate path for HTTPS port")
	fs.String("CONFIG_FILE", "config.yaml", "config file name")
	fs.String("UI_PATH", "./ui", "UI local path")
	fs.String("UI_LOGO", "", "UI logo")
	fs.String("UI_COLOR", "#34577c", "UI color")
	fs.String("UI_MESSAGE", fmt.Sprintf("greetings from service v%v", version.VERSION), "UI message")
	fs.Bool("ENABLE_H2C", false, "allow upgrading to H2C")
	fs.Bool("ENABLE_RANDOM_DELAY", false, "between 0 and 5 seconds random delay by default")
	fs.String("RANDOM_DELAY_UNIT", "s", "either s(seconds) or ms(milliseconds")
	fs.Int("RANDOM_DELAY_MIN_IN_MS", 0, "min for random delay: 0 by default")
	fs.Int("RANDOM_DELAY_MAX_IN_MS", 5, "max for random delay: 5 by default")
	fs.Bool("ENABLE_RANDOM_RANDOM_ERROR", false, "1/3 chances of a random response error")
	fs.Bool("SET_SERVICE_UNHEALTHY", false, "when set, healthy state is never reached")
	fs.Bool("SET_SERVICE_UNREADY", false, "when set, ready state is never reached")
	fs.Bool("ENABLE_CPU_STRESS_TEST", false, "enable cpu stress tests")
	fs.Bool("ENABLE_MEMORY_STRESS_TEST", false, "enable memory stress tests")
	fs.Int("NUMBER_OF_STRESSED_CPU", 0, "number of CPU cores with 100 load")
	fs.Int("DATA_LOADED_IN_MEMORY_FOR_STRESS_TEST_IN_MB", 0, "MB of data to load into memory")
	fs.String("CACHE_SERVER_ADDRESS", "", "Redis address in the format <host>:<port>")

	// authentication service specific flags
	fs.String("AUTHN_USERNAME", "feelguuds", "username of authentication client")
	fs.String("AUTHN_PASSWORD", "feelguuds", "password of authentication client")
	fs.String("AUTHN_ISSUER_BASE_URL", "http://localhost", "authentication service issuer")
	fs.String("AUTHN_ORIGIN", "http://localhost", "origin of auth requests")
	fs.String("AUTHN_DOMAINS", "localhost", "authentication service domains")
	fs.String("AUTHN_PRIVATE_BASE_URL", "http://authentication_service",
		"authentication service private url. should be local host if these are not running on docker containers. "+
			"However if running in docker container with a configured docker network, the url should be equal to the service name")
	fs.String("AUTHN_PUBLIC_BASE_URL", "http://localhost", "authentication service public endpoint")
	fs.String("AUTHN_INTERNAL_PORT", "3000", "authentication service port")
	fs.String("AUTHN_EXTERNAL_PORT", "8000", "authentication service external port")
	fs.Bool("ENABLE_AUTHN_PRIVATE_INTEGRATION", true, "enables communication with authentication service")

	// retry specific configurations
	fs.Int("HTTP_MAX_RETRIES", 5, "max retries to perform on failed http calls")
	fs.Duration("HTTP_MIN_RETRY_WAIT_TIME_IN_MS", 5*time.Millisecond, "minimum time to wait between failed calls for retry")
	fs.Duration("HTTP_MAX_RETRY_WAIT_TIME_IN_MS", 15*time.Millisecond, "maximum time to wait between failed calls for retry")
	fs.Duration("HTTP_REQUEST_TIMEOUT_IN_MS", 300*time.Millisecond, "time until a request is seen as timing out")

	// logging specific configurations
	fs.String("SERVICE_NAME", "authentication_handler_service", "service name")
	fs.String("JAEGER_ENDPOINT", "http://jaeger-collector:14268/api/traces", "jaeger collector endpoint")
	fs.Int("DOWNSTREAM_SERVICE_CONNECTION_LIMIT", 8, "max retries to perform while attempting to connect to downstream services")

	// capture goroutines waiting on synchronization primitives
	runtime.SetBlockProfileRate(1)

	versionFlag := fs.BoolP("ENABLE_VERSION_FROM_FILE", "v", false, "get version number")

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
	viper.RegisterAlias("BACKEND_SERVICE_URLS", "backend-url")
	hostname, _ := os.Hostname()
	viper.SetDefault("JWT_SECRET", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	viper.SetDefault("UI_LOGO", "https://raw.githubusercontent.com/stefanprodan/podinfo/gh-pages/cuddle_clap.gif")
	viper.Set("HOSTNAME", hostname)
	viper.Set("VERSION", version.VERSION)
	viper.Set("REVISION", version.REVISION)
	viper.SetEnvPrefix("SERVICE")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// load config from file
	if _, fileErr := os.Stat(filepath.Join(viper.GetString("CONFIG_PATH"), viper.GetString("CONFIG_FILE"))); fileErr == nil {
		viper.SetConfigName(strings.Split(viper.GetString("CONFIG_FILE"), ".")[0])
		viper.AddConfigPath(viper.GetString("CONFIG_PATH"))
		if readErr := viper.ReadInConfig(); readErr != nil {
			fmt.Printf("Error reading config file, %v\n", readErr)
		}
	} else {
		fmt.Printf("Error to open config file, %v\n", fileErr)
	}

	serviceName := viper.GetString("SERVICE_NAME")
	collectorEndpoint := viper.GetString("JAEGER_ENDPOINT")
	logLevel := viper.GetString("LOG_LEVEL")

	logInstance := core_logging.New(logLevel)
	defer logInstance.ConfigureLogger()
	log := logInstance.Logger

	// initialize metrics object
	coreMetrics := core_metrics.NewCoreMetricsEngineInstance(serviceName, nil)
	serviceMetrics := metrics.New(coreMetrics, serviceName)

	// initialize a tracing object globally
	tracerEngine, closer := tracer.New(serviceName, collectorEndpoint)
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			log.Error("Failed to close handle to tracing instance connection")
		}
	}(closer)

	if tracerEngine == nil {
		log.Fatal("cannot initialize tracer engine")
	}
	opentracing.SetGlobalTracer(tracerEngine.Tracer)

	authnServiceClient := NewAuthServiceClientConnection(err, log)
	if authnServiceClient == nil {
		log.Fatal("failed to initialize connection to authentication service client")
	}

	// start stress tests if any
	numStressedCpus := viper.GetInt("NUMBER_OF_STRESSED_CPU")
	dataInMemForStressTestInMb := viper.GetInt("DATA_LOADED_IN_MEMORY_FOR_STRESS_TEST_IN_MB")
	beginStressTest(numStressedCpus, dataInMemForStressTestInMb, log)

	// validate port
	if _, err := strconv.Atoi(viper.GetString("HTTP_PORT")); err != nil {
		port, _ := fs.GetInt("HTTP_PORT")
		viper.Set("HTTP_PORT", strconv.Itoa(port))
	}

	// validate secure port
	if _, err := strconv.Atoi(viper.GetString("HTTPS_PORT")); err != nil {
		securePort, _ := fs.GetInt("HTTPS_PORT")
		viper.Set("HTTPS_PORT", strconv.Itoa(securePort))
	}

	// validate random delay options
	if viper.GetInt("RANDOM_DELAY_MAX_IN_MS") < viper.GetInt("RANDOM_DELAY_MIN_IN_MS") {
		err := errors.New("`--random-delay-max` should be greater than `--random-delay-min`")
		log.Fatal("please fix configurations", zap.Error(err))
	}

	switch delayUnit := viper.GetString("RANDOM_DELAY_UNIT"); delayUnit {
	case
		"s",
		"ms":
		break
	default:
		err := errors.New("random-delay-unit` accepted values are: s|ms")
		log.Fatal("please fix configurations", zap.Error(err))
	}

	// load gRPC server config
	var grpcCfg grpc.Config
	if err := viper.Unmarshal(&grpcCfg); err != nil {
		err := errors.New("config unmarshal failed")
		log.Fatal("please fix configurations", zap.Error(err))
	}

	// start gRPC server
	if grpcCfg.Port > 0 {
		log.Info("starting grpc server")
		grpcSrv, _ := grpc.NewGRPCServer(&grpcCfg, authnServiceClient, log, serviceMetrics.MicroServiceMetrics, serviceMetrics.Engine, tracerEngine)
		log.Info("successfully started grpc server", zap.Int("port", grpcCfg.Port))
		go grpcSrv.ListenAndServe()
	}

	// load HTTP server config
	var srvCfg api.Config
	if err := viper.Unmarshal(&srvCfg); err != nil {
		log.Fatal("config unmarshal failed", zap.Error(err))
	}

	// log version and port
	log.Info("Starting service",
		zap.String("version", viper.GetString("VERSION")),
		zap.String("revision", viper.GetString("REVISION")),
		zap.String("port", srvCfg.Port),
	)

	// start HTTP server
	srv, _ := api.NewServer(&srvCfg, authnServiceClient, log, serviceMetrics.MicroServiceMetrics, serviceMetrics.Engine, tracerEngine)
	stopCh := signals.SetupSignalHandler()
	srv.ListenAndServe(stopCh)
}

var stressMemoryPayload []byte

// beginStressTest performs cpu and memory stress tests
func beginStressTest(cpus int, mem int, log *zap.Logger) {
	PerformCpuStressTest(cpus, log)
	PerformMemoryStressTest(mem, log)
}

// PerformMemoryStressTest performs memory stress test
func PerformMemoryStressTest(mem int, log *zap.Logger) {
	if mem > 0 {
		path := "/tmp/service.data"
		f, err := os.Create(path)

		if err != nil {
			log.Error("memory stress failed", zap.Error(err))
		}

		if err := f.Truncate(1000000 * int64(mem)); err != nil {
			log.Error("memory stress failed", zap.Error(err))
		}

		stressMemoryPayload, err = ioutil.ReadFile(path)
		f.Close()
		os.Remove(path)
		if err != nil {
			log.Error("memory stress failed", zap.Error(err))
		}
		log.Info("starting CPU stress", zap.Int("memory", len(stressMemoryPayload)))
	}
}

// PerformCpuStressTest performs a cpu stress test
func PerformCpuStressTest(cpus int, log *zap.Logger) {
	done := make(chan int)
	if cpus > 0 {
		log.Info("starting CPU stress", zap.Int("cores", cpus))
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
}

// initializeAuthnClient initializes an instance of the authn client primarily useful in
// communicating with the authentication service securely
func initializeAuthnClient(username, password, audience, issuer, url, origin string) (*core_auth_sdk.Client, error) {
	retryConfig := &core_auth_sdk.RetryConfig{
		MaxRetries:       viper.GetInt("HTTP_MAX_RETRIES"),
		MinRetryWaitTime: viper.GetDuration("HTTP_MIN_RETRY_WAIT_TIME_IN_MS"),
		MaxRetryWaitTime: viper.GetDuration("HTTP_MAX_RETRY_WAIT_TIME_IN_MS"),
		RequestTimeout:   viper.GetDuration("HTTP_REQUEST_TIMEOUT_IN_MS"),
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

// NewAuthServiceClientConnection initializes a new connection to the authentication service and returns a reference to a client object
func NewAuthServiceClientConnection(err error, log *zap.Logger) *core_auth_sdk.Client {
	// initialize authentication client in order to establish communication with the
	// authentication service. This serves as a singular source of truth for authentication needs
	authUsername := viper.GetString("AUTHN_USERNAME")
	authPassword := viper.GetString("AUTHN_PASSWORD")
	domains := viper.GetString("AUTHN_DOMAINS")
	privateURL := viper.GetString("AUTHN_PRIVATE_BASE_URL") + ":" + viper.GetString("AUTHN_INTERNAL_PORT")
	origin := viper.GetString("AUTHN_ORIGIN")
	issuer := viper.GetString("AUTHN_ISSUER_BASE_URL") + ":" + viper.GetString("AUTHN_EXTERNAL_PORT")

	authnClient, err := initializeAuthnClient(authUsername, authPassword, domains, issuer, privateURL, origin)
	// crash the process if we cannot connect to the authentication service
	if err != nil {
		log.Fatal("failed to initialized authentication service client", zap.Error(err))
	}

	if err = connectToAuthenticationService(authnClient, log); err != nil {
		log.Fatal("failed to initiate connection to downstream service", zap.Error(err))
	}

	return authnClient
}

// connectToAuthenticationService attempts to connect to a downstream authentication service
func connectToAuthenticationService(authnClient *core_auth_sdk.Client, log *zap.Logger) error {
	retryLimit := viper.GetInt("DOWNSTREAM_SERVICE_CONNECTION_LIMIT")

	var response = make(chan interface{}, 1)
	err := retry.Do(
		func(conn chan<- interface{}) func() error {
			return func() error {
				opResponse, err := authnClient.ServerStats()
				if err != nil {
					return err
				}
				response <- opResponse
				return nil
			}
		}(response),
		retry.MaxTries(retryLimit),
		retry.Timeout(time.Millisecond*time.Duration(500)),
		retry.Sleep(time.Millisecond*time.Duration(50)),
	)
	return err
}
