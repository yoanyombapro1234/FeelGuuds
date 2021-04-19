package grpc

import (
	"context"
	"time"

	"github.com/giantswarm/retry-go"
	"github.com/opentracing/opentracing-go"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

type RpcOperationFunc func() (interface{}, error)

// setCtxRequestTimeout sets the request deadline in the context. This function should be invoked prior to any rpc calls
func (s *Server) setCtxRequestTimeout(ctx context.Context) context.Context {
	clientDeadline := time.Now().Add(time.Duration(s.config.RpcDeadline) * time.Millisecond)
	ctx, _ = context.WithDeadline(ctx, clientDeadline)
	return ctx
}

// performRetryableRpcCall performs an rpc call using retries in the face of errors
func (s *Server) performRetryableRpcCall(ctx context.Context, f func() (interface{}, error)) (interface{}, error) {
	var response = make(chan *interface{}, 1)

	err := retry.Do(
		func(conn chan<- *interface{}) func() error {
			return func() error {
				opResponse, err := f()
				if err != nil {
					return err
				}
				response <- &opResponse
				return nil
			}
		}(response),
		retry.MaxTries(s.config.RpcRetries),
		retry.Timeout(time.Millisecond*time.Duration(s.config.RpcDeadline)),
		retry.Sleep(time.Millisecond*time.Duration(s.config.RpcRetryBackoff)),
	)

	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.Canceled {
		return nil, service_errors.ErrRequestTimeout
	}

	return <-response, nil
}

// PerformRetryableRPCOperation performs a retryable operation
func (s *Server) PerformRetryableRPCOperation(ctx context.Context, span opentracing.Span, op RpcOperationFunc, opType string) RpcOperationFunc {
	return func() (interface{}, error) {
		var (
			begin = time.Now()
			took  = time.Since(begin)
		)

		retryableOp := func() (interface{}, error) {
			return s.performRetryableRpcCall(ctx, op)
		}

		ctx = opentracing.ContextWithSpan(ctx, span)
		return s.performRPCOperationAndInstrument(ctx, retryableOp, opType, &took)
	}
}
