package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	core_metrics "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-metrics"
)

type MetricsEngine struct {
	MicroServiceMetrics *CoreMetrics
	Engine              *core_metrics.CoreMetricsEngine
}

func NewMetricsEngine(engine *core_metrics.CoreMetricsEngine, serviceName string) *MetricsEngine {
	return &MetricsEngine{
		MicroServiceMetrics: NewCoreMetrics(engine, serviceName),
		Engine:              engine,
	}
}

type CoreMetrics struct {
	ServiceName string
	// tracks the number of http requests partitioned by name and status code
	// used for monitoring and alerting (RED method)
	HttpRequestCounter *core_metrics.CounterVec
	// tracks the latencies associated with a http requests by operation name
	// used for horizontal pod auto-scaling (Kubernetes HPA v2)
	HttpRequestLatencyCounter *core_metrics.HistogramVec
	// tracks the number of times there was a failure or success when trying to extract id from the request url
	ExtractIdOperationCounter *core_metrics.CounterVec
	// tracks the number of times there was a failure or success when trying to extract id from the request url
	RemoteOperationStatusCounter    *core_metrics.CounterVec
	RemoteOperationsLatencyCounter  *core_metrics.HistogramVec
	InvalidRequestParametersCounter *core_metrics.CounterVec
	CastingOperationFailureCounter  *core_metrics.CounterVec
	DecodeRequestStatusCounter      *core_metrics.CounterVec
}

func NewCoreMetrics(engine *core_metrics.CoreMetricsEngine, serviceName string) *CoreMetrics {
	return &CoreMetrics{
		ServiceName:                     serviceName,
		HttpRequestCounter:              NewHttpRequestCounter(engine, serviceName),
		HttpRequestLatencyCounter:       NewHttpRequestLatencyCounter(engine, serviceName),
		ExtractIdOperationCounter:       NewExtractIdOperationCounter(engine, serviceName),
		RemoteOperationStatusCounter:    NewRemoteOperationStatusCounter(engine, serviceName),
		RemoteOperationsLatencyCounter:  NewRemoteOperationLatencyCounter(engine, serviceName),
		InvalidRequestParametersCounter: NewInvalidRequestParametersCounter(engine, serviceName),
		CastingOperationFailureCounter:  NewCastingOperationFailureCounter(engine, serviceName),
		DecodeRequestStatusCounter:      NewDecodeRequestStatusCounter(engine, serviceName),
	}
}

func NewHttpRequestCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.CounterVec {
	newCounter := core_metrics.NewCounterVec(&core_metrics.CounterOpts{
		Namespace: serviceName,
		Subsystem: "HTTP",
		Name:      fmt.Sprintf("%s_http_requests_total", serviceName),
		Help:      "How many HTTP requests processed partitioned by name and status code",
	}, []string{"name", "code"})

	engine.RegisterMetric(newCounter)
	return newCounter
}

func NewHttpRequestLatencyCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.HistogramVec {
	newCounter := core_metrics.NewHistogramVec(&core_metrics.HistogramOpts{
		Namespace:         serviceName,
		Subsystem:         "HTTP",
		Name:              fmt.Sprintf("%s_http_requests_latencies", serviceName),
		Help:              "Seconds spent serving HTTP requests.",
		ConstLabels:       nil,
		Buckets:           prometheus.DefBuckets,
		DeprecatedVersion: "",
		StabilityLevel:    "",
	}, []string{"method", "path", "status"})
	engine.RegisterMetric(newCounter)
	return newCounter
}

func NewExtractIdOperationCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.CounterVec {
	// tracks the number of times there was a failure or success when trying to extract id from the request url
	newCounter := core_metrics.NewCounterVec(&core_metrics.CounterOpts{
		Namespace: serviceName,
		Subsystem: "HTTP",
		Name:      fmt.Sprintf("%s_status_of_extract_id_operation_from_requests_total", serviceName),
		Help:      "The status of the extract the id operation from the HTTP requests processed partitioned by operation name and operation status",
	}, []string{"operation_name", "status"})
	engine.RegisterMetric(newCounter)
	return newCounter
}

func NewRemoteOperationStatusCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.CounterVec {
	newCounter := core_metrics.NewCounterVec(&core_metrics.CounterOpts{
		Namespace: serviceName,
		Subsystem: "HTTP",
		Name:      fmt.Sprintf("%s_status_of_remote_operation_total", serviceName),
		Help:      "A count of the status all remote operations operation",
	}, []string{"operation_name", "status"})
	engine.RegisterMetric(newCounter)
	return newCounter
}

func NewRemoteOperationLatencyCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.HistogramVec {
	newCounter := core_metrics.NewHistogramVec(&core_metrics.HistogramOpts{
		Namespace:         serviceName,
		Subsystem:         "HTTP",
		Name:              fmt.Sprintf("%s_remote_operation_requests_latencies", serviceName),
		Help:              "Seconds spent serving remote operations HTTP requests.",
		ConstLabels:       nil,
		Buckets:           prometheus.DefBuckets,
		DeprecatedVersion: "",
		StabilityLevel:    "",
	}, []string{"operation", "status"})
	engine.RegisterMetric(newCounter)
	return newCounter
}

func NewInvalidRequestParametersCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.CounterVec {
	newCounter := core_metrics.NewCounterVec(&core_metrics.CounterOpts{
		Namespace: serviceName,
		Subsystem: "HTTP",
		Name:      fmt.Sprintf("%s_invalid_request_parameters_total", serviceName),
		Help:      "A count of the total number of invalid request parameter count",
	}, []string{"operation_name"})
	engine.RegisterMetric(newCounter)
	return newCounter
}

func NewCastingOperationFailureCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.CounterVec {
	newCounter := core_metrics.NewCounterVec(&core_metrics.CounterOpts{
		Namespace: serviceName,
		Subsystem: "HTTP",
		Name:      fmt.Sprintf("%s_casting_operation_failure_total", serviceName),
		Help:      "A count of the total number of failed casts from interface to object",
	}, []string{"operation_name"})
	engine.RegisterMetric(newCounter)
	return newCounter
}

func NewDecodeRequestStatusCounter(engine *core_metrics.CoreMetricsEngine, serviceName string) *core_metrics.CounterVec {
	newCounter := core_metrics.NewCounterVec(&core_metrics.CounterOpts{
		Namespace: serviceName,
		Subsystem: "HTTP",
		Name:      fmt.Sprintf("%s_decoder_request_op_counter_total", serviceName),
		Help:      "A count of the status of all decode operations",
	}, []string{"operation_name", "status"})
	engine.RegisterMetric(newCounter)
	return newCounter
}
