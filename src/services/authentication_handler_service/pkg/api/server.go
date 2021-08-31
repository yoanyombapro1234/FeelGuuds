package api

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/version"
	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-auth-sdk"
	core_metrics "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-metrics"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds_Core/core/core-tracing/jaeger"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	_ "github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/api/docs"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/fscache"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/metrics"
)

// @title Service API
// @version 2.0
// @description Go microservice template for Kubernetes.

// @contact.name Source Code
// @contact.url https://github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service

// @license.name MIT License
// @license.url https://github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/blob/master/LICENSE

// @host localhost:9898
// @BasePath /
// @schemes http https

var (
	healthy int32
	ready   int32
	watcher *fscache.Watcher
)

type Config struct {
	HttpClientTimeout         time.Duration `mapstructure:"HTTP_CLIENT_TIMEOUT_IN_MINUTES"`
	HttpServerTimeout         time.Duration `mapstructure:"HTTP_SERVER_TIMEOUT_IN_SECONDS"`
	HttpServerShutdownTimeout time.Duration `mapstructure:"HTTP_SERVER_SHUTDOWN_TIMEOUT_IN_SECONDS"`
	BackendURL                []string      `mapstructure:"BACKEND_SERVICE_URLS"`
	UILogo                    string        `mapstructure:"UI_LOGO"`
	UIMessage                 string        `mapstructure:"UI_MESSAGE"`
	UIColor                   string        `mapstructure:"UI_COLOR"`
	UIPath                    string        `mapstructure:"UI_PATH"`
	DataPath                  string        `mapstructure:"DATA_PATH"`
	ConfigPath                string        `mapstructure:"CONFIG_PATH"`
	CertPath                  string        `mapstructure:"CERT_PATH"`
	Port                      string        `mapstructure:"HTTP_PORT"`
	SecurePort                string        `mapstructure:"HTTPS_PORT"`
	PortMetrics               int           `mapstructure:"METRICS_PORT"`
	Hostname                  string        `mapstructure:"HOSTNAME"`
	H2C                       bool          `mapstructure:"H2C"`
	RandomDelay               bool          `mapstructure:"ENABLE_RANDOM_DELAY"`
	RandomDelayUnit           string        `mapstructure:"RANDOM_DELAY_UNIT"`
	RandomDelayMin            int           `mapstructure:"RANDOM_DELAY_MIN_IN_MS"`
	RandomDelayMax            int           `mapstructure:"RANDOM_DELAY_MAX_IN_MS"`
	RandomError               bool          `mapstructure:"ENABLE_RANDOM_ERROR"`
	Unhealthy                 bool          `mapstructure:"SET_SERVICE_UNHEALTHY"`
	Unready                   bool          `mapstructure:"SET_SERVICE_UNREADY"`
	JWTSecret                 string        `mapstructure:"JWT_SECRET"`
	CacheServer               string        `mapstructure:"CACHE_SERVER_ADDRESS"`
}

type Server struct {
	router        *mux.Router
	config        *Config
	pool          *redis.Pool
	handler       http.Handler
	authnClient   *core_auth_sdk.Client
	logger        *zap.Logger
	metrics       *metrics.CoreMetrics
	metricsEngine *core_metrics.CoreMetricsEngine
	tracerEngine  *core_tracing.TracingEngine
}

func NewServer(config *Config, client *core_auth_sdk.Client, logger *zap.Logger, metrics *metrics.CoreMetrics,
	metricsEngine *core_metrics.CoreMetricsEngine, tracer *core_tracing.TracingEngine) (*Server, error) {
	srv := &Server{
		router:        mux.NewRouter(),
		config:        config,
		pool:          nil,
		handler:       nil,
		authnClient:   client,
		logger:        logger,
		metrics:       metrics,
		metricsEngine: metricsEngine,
		tracerEngine:  tracer,
	}

	return srv, nil
}

func (s *Server) registerHandlers() {
	s.router.Handle("/metrics", promhttp.Handler())
	s.router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	s.router.HandleFunc("/", s.indexHandler).HeadersRegexp("User-Agent", "^Mozilla.*").Methods("GET")
	s.router.HandleFunc("/healthz", s.healthzHandler).Methods("GET")
	s.router.HandleFunc("/readyz", s.readyzHandler).Methods("GET")
	s.router.HandleFunc("/readyz/enable", s.enableReadyHandler).Methods("POST")
	s.router.HandleFunc("/readyz/disable", s.disableReadyHandler).Methods("POST")
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	s.router.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		doc, err := swag.ReadDoc()
		if err != nil {
			s.logger.Error("swagger error", zap.Error(err), zap.String("path", "/swagger.json"))
		}
		w.Write([]byte(doc))
	})
}

func (s *Server) registerMiddlewares() {
	s.router.Use(versionMiddleware)
	if s.config.RandomError {
		s.router.Use(randomErrorMiddleware)
	}
}

func (s *Server) ListenAndServe(stopCh <-chan struct{}) {
	go s.startMetricsServer()

	s.registerHandlers()
	s.registerMiddlewares()

	if s.config.H2C {
		s.handler = h2c.NewHandler(s.router, &http2.Server{})
	} else {
		s.handler = s.router
	}

	// s.printRoutes()

	// load configs in memory and start watching for changes in the config dir
	if stat, err := os.Stat(s.config.ConfigPath); err == nil && stat.IsDir() {
		var err error
		watcher, err = fscache.NewWatch(s.config.ConfigPath)
		if err != nil {
			s.logger.Error("config watch error", zap.Error(err), zap.String("path", s.config.ConfigPath))
		} else {
			watcher.Watch()
		}
	}

	// start redis connection pool
	ticker := time.NewTicker(30 * time.Second)
	s.startCachePool(ticker, stopCh)

	// create the http server
	srv := s.startServer()

	// create the secure server
	secureSrv := s.startSecureServer()

	// signal Kubernetes the server is ready to receive traffic
	if !s.config.Unhealthy {
		atomic.StoreInt32(&healthy, 1)
	}
	if !s.config.Unready {
		atomic.StoreInt32(&ready, 1)
	}

	// wait for SIGTERM or SIGINT
	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), s.config.HttpServerShutdownTimeout)
	defer cancel()

	// all calls to /healthz and /readyz will fail from now on
	atomic.StoreInt32(&healthy, 0)
	atomic.StoreInt32(&ready, 0)

	// close cache pool
	if s.pool != nil {
		_ = s.pool.Close()
	}

	s.logger.Info("Shutting down HTTP/HTTPS server", zap.Duration("timeout", s.config.HttpServerShutdownTimeout))

	// wait for Kubernetes readiness probe to remove this instance from the load balancer
	// the readiness check interval must be lower than the timeout
	if viper.GetString("level") != "debug" {
		time.Sleep(3 * time.Second)
	}

	// determine if the http server was started
	if srv != nil {
		if err := srv.Shutdown(ctx); err != nil {
			s.logger.Error("HTTP server graceful shutdown failed", zap.Error(err))
		}
	}

	// determine if the secure server was started
	if secureSrv != nil {
		if err := secureSrv.Shutdown(ctx); err != nil {
			s.logger.Error("HTTPS server graceful shutdown failed", zap.Error(err))
		}
	}
}

func (s *Server) startServer() *http.Server {

	// determine if the port is specified
	if s.config.Port == "0" {

		// move on immediately
		return nil
	}

	srv := &http.Server{
		Addr:         ":" + s.config.Port,
		WriteTimeout: s.config.HttpServerTimeout,
		ReadTimeout:  s.config.HttpServerTimeout,
		IdleTimeout:  2 * s.config.HttpServerTimeout,
		Handler:      s.handler,
	}

	// start the server in the background
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatal("HTTP server crashed", zap.Error(err))
		}
	}()

	// return the server and routine
	return srv
}

func (s *Server) startSecureServer() *http.Server {

	// determine if the port is specified
	if s.config.SecurePort == "0" {

		// move on immediately
		return nil
	}

	srv := &http.Server{
		Addr:         ":" + s.config.SecurePort,
		WriteTimeout: s.config.HttpServerTimeout,
		ReadTimeout:  s.config.HttpServerTimeout,
		IdleTimeout:  2 * s.config.HttpServerTimeout,
		Handler:      s.handler,
	}

	cert := path.Join(s.config.CertPath, "/cert.pem")
	key := path.Join(s.config.CertPath, "/key.unencrypted.pem")

	// start the server in the background
	go func() {
		if err := srv.ListenAndServeTLS(cert, key); err != http.ErrServerClosed {
			s.logger.Fatal("HTTPS server crashed", zap.Error(err))
		}
	}()

	// return the server
	return srv
}

func (s *Server) startMetricsServer() {
	if s.config.PortMetrics > 0 {
		mux := http.DefaultServeMux
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%v", s.config.PortMetrics),
			Handler: mux,
		}

		srv.ListenAndServe()
	}
}

func (s *Server) printRoutes() {
	s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})
}

func (s *Server) startCachePool(ticker *time.Ticker, stopCh <-chan struct{}) {
	if s.config.CacheServer == "" {
		return
	}
	s.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", s.config.CacheServer)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	// set <hostname>=<version> with an expiry time of one minute
	setVersion := func() {
		conn := s.pool.Get()
		if _, err := conn.Do("SET", s.config.Hostname, version.VERSION, "EX", 60); err != nil {
			s.logger.Warn("cache server is offline", zap.Error(err), zap.String("server", s.config.CacheServer))
		}
		_ = conn.Close()
	}

	// set version on a schedule
	go func() {
		setVersion()
		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				setVersion()
			}
		}
	}()
}

type ArrayResponse []string
type MapResponse map[string]string
