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

// GetAccount obtains an account as long as the account exists from the context of the authentication service
func (s *Server) GetAccount(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
	ctx = s.setCtxRequestTimeout(ctx)
	ctx, parentSpan := s.StartRootSpan(ctx, constants.GET_ACCOUNT)
	defer parentSpan.Finish()

	if req == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	if req.Id == 0 {
		s.metrics.InvalidRequestParametersCounter.WithLabelValues(constants.GET_ACCOUNT).Inc()
		err := service_errors.ErrInvalidInputArguments
		var msg = "invalid input parameters. please specify a valid user id"
		s.logger.Error(err, msg)
		return nil, err
	}

	var (
		begin           = time.Now()
		took            = time.Since(begin)
		operation       = func() (interface{}, error) {
			account, err := s.authnClient.GetAccount(strconv.Itoa(int(req.GetId())))
			if err != nil {
				return nil, err
			}
			return account, nil
		}
		retryableOperation = func() (interface{}, error) {
			return s.performRetryableRpcCall(ctx, operation)
		}
	)

	ctx = opentracing.ContextWithSpan(ctx, parentSpan)
	result, err := s.PerformRPCOperationAndInstrument(ctx, retryableOperation, constants.DELETE_ACCOUNT, &took)
	if err != nil {
		s.logger.Error(err, err.Error())
		return nil, err
	}

	account, ok := result.(*proto.Account)
	if !ok {
		s.metrics.CastingOperationFailureCounter.WithLabelValues(constants.GET_ACCOUNT)

		err := service_errors.ErrFailedToCastAccount
		s.logger.For(ctx).ErrorM(err, err.Error())
		return nil, err
	}

	s.logger.For(ctx).Info("Successfully obtained user account", zap.Int("Id", int(req.GetId())))
	return &proto.GetAccountResponse{
		Account:              account,
		Error:                "",
	}, nil
}
