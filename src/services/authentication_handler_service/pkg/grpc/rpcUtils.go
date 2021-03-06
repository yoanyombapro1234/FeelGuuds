package grpc

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	var response = make(chan interface{}, 1)

	/*
		err := retry.Do(
			func(conn chan<- interface{}) func() error {
				return func() error {
					opResponse, err := f()
					if err != nil {
						return err
					}
					response <- opResponse
					return nil
				}
			}(response),
			retry.MaxTries(s.config.RpcRetries),
			retry.Timeout(time.Millisecond*time.Duration(s.config.RpcDeadline)),
			retry.Sleep(time.Millisecond*time.Duration(s.config.RpcRetryBackoff)),
		)

		retries := 1
		for retries < 4 {
			// perform a test request to the authentication service
			data, err := client.ServerStats()
			if err != nil {
				if retries != 4 {
					logger.Error(err, "failed to connect to authentication service")
				} else {
					logger.Fatal(err, "failed to connect to authentication service")
				}
				retries += 1
			} else {
				retries = 4
				logger.Info("data", zap.Any("result", data))
			}

			time.Sleep(1 * time.Second)
		}
	*/

	err := func(conn chan<- interface{}) func() error {
		return func() error {
			opResponse, err := f()
			if err != nil {
				return err
			}
			response <- opResponse
			return nil
		}
	}(response)()

	if err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	if ctx.Err() == context.Canceled {
		return nil, service_errors.ErrRequestTimeout
	}

	return <-response, nil
}

// PerformRetryableRPCOperation performs a retryable operation
func (s *Server) PerformRetryableRPCOperation(ctx context.Context, span opentracing.Span, op proto.DownStreamOperation, opType string) proto.DownStreamOperation {
	return func() (interface{}, error) {
		var (
			begin = time.Now()
			took  = time.Since(begin)
		)

		ctx = opentracing.ContextWithSpan(ctx, span)

		retryableOp := func() (interface{}, error) {
			s.logger.For(ctx).Info("performing retryable http operation", "operation type", opType)
			return s.performRetryableRpcCall(ctx, op)
		}

		ctx = opentracing.ContextWithSpan(ctx, span)
		return s.performRPCOperationAndInstrument(ctx, retryableOp, opType, &took)
	}
}

// ConfigureAndStartRootSpan configures a parent span object and starts it
func (s *Server) ConfigureAndStartRootSpan(ctx context.Context, operationType string) (context.Context, opentracing.Span) {
	ctx = s.setCtxRequestTimeout(ctx)
	ctx, rootSpan := s.StartRootSpan(ctx, operationType)
	return ctx, rootSpan
}
