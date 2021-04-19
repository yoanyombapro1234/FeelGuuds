package grpc

import (
	"context"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	core_auth_sdk "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-auth-sdk"
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
		s.logger.Error(err, "invalid input parameters. please specify a valid user id")
		return nil, err
	}

	var (
		operation = func() (interface{}, error) {
			account, err := s.authnClient.GetAccount(strconv.Itoa(int(req.GetId())))
			if err != nil {
				return nil, err
			}
			return account, nil
		}
	)

	ctx = opentracing.ContextWithSpan(ctx, parentSpan)
	result, err := s.PerformRetryableRPCOperation(ctx, parentSpan, operation, constants.GET_ACCOUNT)()
	if err != nil {
		s.logger.Error(err, err.Error())
		return nil, err
	}

	account, ok := result.(*core_auth_sdk.Account)
	if !ok {
		s.metrics.CastingOperationFailureCounter.WithLabelValues(constants.GET_ACCOUNT)

		err := service_errors.ErrFailedToCastAccount
		s.logger.For(ctx).ErrorM(err, err.Error())
		return nil, err
	}

	s.logger.For(ctx).Info("Successfully obtained user account", zap.Int("Id", int(req.GetId())))
	return &proto.GetAccountResponse{
		Account: &proto.Account{
			Id:                   uint32(account.ID),
			Username:             account.Username,
			Locked:               account.Locked,
			Deleted:              account.Deleted,
			XXX_sizecache:        0,
		},
		Error:   "",
	}, nil
}
