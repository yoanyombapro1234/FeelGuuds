package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
)

// PerformRPCOperationAndInstrument performs a rpc call to an external service and instruments the resulting response accordingly
func (s *Server) PerformRPCOperationAndInstrument(
	ctx context.Context,
	f func() (interface{}, error),
	operationType string,
	took *time.Duration) (interface{}, error) {

	authnSvcRpcSpan := s.tracerEngine.CreateChildSpan(ctx, fmt.Sprintf("AUTHENTICATION_SERVICE_%s_RPC_REQUEST", operationType))
	defer authnSvcRpcSpan.Finish()

	var status = constants.SUCCESS
	result, err := f()
	// error scenarios are only when the rpc request took too long or an error occured
	if err != nil || ctx.Err() == context.Canceled {
		status = constants.FAILURE
	}

	s.metrics.RemoteOperationStatusCounter.WithLabelValues(operationType, status).Inc()
	s.metrics.RemoteOperationsLatencyCounter.WithLabelValues(operationType, status).Observe(took.Seconds())
	return result, err
}

// StartRootSpan starts the rootspan of the current operation at hand
func (s *Server) StartRootSpan(ctx context.Context, operationType string) (context.Context, opentracing.Span) {
	s.logger.For(ctx).InfoM("GRPC request received", zap.String("method", operationType))

	spanCtx, _ := s.tracerEngine.Tracer.Extract(opentracing.HTTPHeaders, nil)
	parentSpan := s.tracerEngine.Tracer.StartSpan(operationType, ext.RPCServerOption(spanCtx))
	ctx = opentracing.ContextWithSpan(ctx, parentSpan)

	return ctx, parentSpan
}
