package grpc_test_utils

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-tracing"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
)

// InitializeTracingEngine initiaize a tracing object globally
func InitializeTracingEngine(serviceName string) (*core_tracing.TracingEngine, io.Closer) {
	return core_tracing.NewTracer(serviceName, constants.COLLECTOR_ENDPOINT, prometheus.New())
}

// InitializeLoggingEngine initializes logging object
func InitializeLoggingEngine(ctx context.Context) core_logging.ILog {
	rootSpan := opentracing.SpanFromContext(ctx)
	logger := core_logging.NewJSONLogger(nil, rootSpan)
	return logger
}