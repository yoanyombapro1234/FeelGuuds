package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/database"
	svcErr "github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/service_errors"
	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/api"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/grpc"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/signals"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/version"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-tracing"
)

func main() {
	// flags definition
	fs, err := DefineAndParseFlags()

	// load config from file
	LoadConfigurationsFromFile()

	// configure distributed tracing
	tracingEngine, closer := SetupDistributedTracingEngine()
	defer closer.Close()

	if tracingEngine == nil {
		panic("cannot initialize tracer engine")
	}
	opentracing.SetGlobalTracer(tracingEngine.Tracer)

	// start root span
	ctx := context.Background()
	rootSpan := opentracing.SpanFromContext(ctx)

	// configure logging
	logger := core_logging.NewJSONLogger(nil, rootSpan)

	// configure db connection
	db, err := ConnectToDatabase(ctx, logger, tracingEngine)
	if err != nil {
		logger.For(ctx).Fatal(err, err.Error())
	}

	// start stress tests if any
	beginStressTest(viper.GetInt("STRESS_CPU"), viper.GetInt("STRESS_MEMORY"), logger)
	// validate http and https ports
	ValidateHttpPorts(fs)
	// validate random delay options
	ValidateRandomDelayOptions(logger)
	// load gRPC server config
	grpcCfg := LoadGrpcConfigs(logger)

	// start gRPC server
	if grpcCfg.Port > 0 {
		grpcSrv, _ := grpc.NewServer(&grpc.ServerInitializationParams{
			Config:             &grpcCfg,
			Logger:             logger,
			TracerEngine:       tracingEngine,
			DatabaseConnection: db,
		})

		go grpcSrv.ListenAndServe(false)
	}

	srvCfg := LoadHttpConfigs(logger)

	// log version and port
	logger.Info("Starting service",
		zap.String("VERSION", viper.GetString("VERSION")),
		zap.String("REVISION", viper.GetString("REVISION")),
		zap.String("HTTP_PORT", srvCfg.Port),
	)

	// start HTTP server
	srv, _ := api.NewServer(&srvCfg, logger)
	stopCh := signals.SetupSignalHandler()
	srv.ListenAndServe(stopCh)
}

func DefineAndParseFlags() (*pflag.FlagSet, error) {
	fs := DefineFlags()
	versionFlag := fs.BoolP("VERSION", "v", false, "get version number")
	// parse flags
	err := ParseFlags(fs, versionFlag)

	// bind flags and environment variables
	BindFlagsToEnvironmentVariables(fs)
	return fs, err
}

func DefineFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("default", pflag.ContinueOnError)
	fs.Int("HTTP_PORT", 9898, "HTTP PORT")
	fs.Int("secure-port", 0, "HTTPS port")
	fs.Int("port-metrics", 0, "metrics port")
	fs.Int("grpc-port", 0, "gRPC port")
	fs.String("grpc-service-name", "service", "gPRC service name")
	fs.String("level", "info", "log level debug, info, warn, error, flat or panic")
	fs.StringSlice("BACKEND_URL", []string{}, "backend service URL")
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
	// TODO: reconfigure this to leverage datadog instead
	fs.String("JAEGER_ENDPOINT", "http://jaeger-collector:14268/api/traces", "jaeger collector endpoint")
	return fs
}

func LoadHttpConfigs(logger core_logging.ILog) api.Config {
	// load HTTP server config
	var srvCfg api.Config
	if err := viper.Unmarshal(&srvCfg); err != nil {
		logger.Fatal(svcErr.ErrInvalidEnvironmentVariableConfigurations, "config unmarshal failed")
	}
	return srvCfg
}

func LoadGrpcConfigs(logger core_logging.ILog) grpc.Config {
	var grpcCfg grpc.Config
	if err := viper.Unmarshal(&grpcCfg); err != nil {
		logger.Fatal(svcErr.ErrInvalidEnvironmentVariableConfigurations, "config unmarshal failed")
	}
	return grpcCfg
}

func ValidateRandomDelayOptions(logger core_logging.ILog) {
	// validate random delay options
	if viper.GetInt("RANDOM_DELAY_MAX") < viper.GetInt("RANDOM_DELAY_MIN") {
		logger.Fatal(svcErr.ErrInvalidEnvironmentVariableConfigurations, "`--RANDOM_DELAY_MAX` should be greater than `--RANDOM_DELAY_MIN`")
	}

	switch delayUnit := viper.GetString("RANDOM_DELAY_UNIT"); delayUnit {
	case
		"s",
		"ms":
		break
	default:
		logger.Fatal(svcErr.ErrInvalidEnvironmentVariableConfigurations, "`random-delay-unit` accepted values are: s|ms")
	}
}

func ValidateHttpPorts(fs *pflag.FlagSet) {
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
}

func ParseFlags(fs *pflag.FlagSet, versionFlag *bool) error {
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
	return err
}

func SetupDistributedTracingEngine() (*core_tracing.TracingEngine, io.Closer) {
	svcName := viper.GetString("GRPC_SERVICE_NAME")
	collectorEndpoint := viper.GetString("JAEGER_ENDPOINT")
	tracingEngine, closer := core_tracing.NewTracer(svcName, collectorEndpoint, prometheus.New())
	return tracingEngine, closer
}

func LoadConfigurationsFromFile() {
	if _, fileErr := os.Stat(filepath.Join(viper.GetString("CONFIG_PATH"), viper.GetString("CONFIG"))); fileErr == nil {
		viper.SetConfigName(strings.Split(viper.GetString("CONFIG"), ".")[0])
		viper.AddConfigPath(viper.GetString("CONFIG_PATH"))
		if readErr := viper.ReadInConfig(); readErr != nil {
			fmt.Printf("Error reading config file, %v\n", readErr)
		}
	} else {
		fmt.Printf("Error to open config file, %v\n", fileErr)
	}
}

func BindFlagsToEnvironmentVariables(fs *pflag.FlagSet) {
	viper.BindPFlags(fs)
	viper.RegisterAlias("backendUrl", "BACKEND_URL")
	hostname, _ := os.Hostname()
	viper.SetDefault("JWT_SECRET", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	viper.SetDefault("UI_LOGO", "https://raw.githubusercontent.com/stefanprodan/podinfo/gh-pages/cuddle_clap.gif")
	viper.Set("HOSTNAME", hostname)
	viper.Set("VERSION", version.VERSION)
	viper.Set("REVISION", version.REVISION)
	viper.SetEnvPrefix("SERVICE")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func ConnectToDatabase(ctx context.Context, logger core_logging.ILog, tracingEngine *core_tracing.TracingEngine) (*database.Db, error) {
	connectionString := GetDatabaseConnectionString()
	logger.Info(fmt.Sprintf("Database connection string : %s ", connectionString))

	return database.New(ctx, database.ConnectionInitializationParams{
		ConnectionString:       connectionString,
		TracingEngine:          tracingEngine,
		Logger:                 logger,
		MaxConnectionAttempts:  4,
		MaxRetriesPerOperation: 4,
		RetryTimeOut:           3 * time.Second,
		RetrySleepInterval:     50 * time.Millisecond,
	})
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
			logger.Error(err, "memory stress failed")
		}

		if err := f.Truncate(1000000 * int64(mem)); err != nil {
			logger.Error(err, "memory stress failed")
		}

		stressMemoryPayload, err = ioutil.ReadFile(path)
		f.Close()
		os.Remove(path)
		if err != nil {
			logger.Error(err, "memory stress failed")
		}
		logger.Info("starting CPU stress", zap.Int("memory", len(stressMemoryPayload)))
	}
}

func GetDatabaseConnectionString() string {
	host := viper.GetString("db_host")
	port := viper.GetInt("db_port")
	user := viper.GetString("db_user")
	password := viper.GetString("db_password")
	dbname := viper.GetString("db_name")
	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return connectionString
}
