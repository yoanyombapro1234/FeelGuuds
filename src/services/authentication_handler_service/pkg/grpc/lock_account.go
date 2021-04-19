package grpc

import (
	"context"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

// LockAccount locks an account as long as it exists from the context of the authentication service
func (s *Server) LockAccount(ctx context.Context, req *proto.LockAccountRequest) (*proto.LockAccountResponse, error) {
	ctx = s.setCtxRequestTimeout(ctx)
	ctx, parentSpan := s.StartRootSpan(ctx, constants.LOCK_ACCOUNT)
	defer parentSpan.Finish()

	if req == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	if req.Id == 0 {
		s.metrics.InvalidRequestParametersCounter.WithLabelValues(constants.LOCK_ACCOUNT).Inc()
		err := service_errors.ErrInvalidInputArguments
		var msg = "invalid input parameters. please specify a valid user id"
		s.logger.Error(err, msg)
		return nil, err
	}

	var (
		begin           = time.Now()
		took            = time.Since(begin)
		operation       = func() (interface{}, error) {
			if err := s.authnClient.LockAccount(strconv.Itoa(int(req.GetId()))); err != nil {
				return nil, err
			}
			return nil, nil
		}
		retryableOperation = func() (interface{}, error) {
			return s.performRetryableRpcCall(ctx, operation)
		}
	)

	ctx = opentracing.ContextWithSpan(ctx, parentSpan)
	_, err := s.PerformRPCOperationAndInstrument(ctx, retryableOperation, constants.LOCK_ACCOUNT, &took)
	if err != nil {
		s.logger.Error(err, err.Error())
		return nil, err
	}

	s.logger.For(ctx).Info("Successfully locked user account", zap.Int("Id", int(req.GetId())))
	return &proto.LockAccountResponse{
		Error:                "",
	}, nil
}
