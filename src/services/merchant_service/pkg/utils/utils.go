package utils

import (
	"context"

	"github.com/opentracing/opentracing-go"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
)

func StartRootOperationSpan(ctx context.Context, operationType string, logger core_logging.ILog) (context.Context,
	opentracing.Span) {
	logger.For(ctx).Info("Starting parent span for operation")
	span, ctx := opentracing.StartSpanFromContext(ctx, operationType)
	ctx = opentracing.ContextWithSpan(ctx, span)
	return ctx, span
}
