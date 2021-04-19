package grpc

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

// LogoutAccount revokes the user account session from the context of the authentication service
func (s *Server) LogoutAccount(ctx context.Context, req *proto.LogoutAccountRequest) (*proto.LogoutAccountResponse, error) {
	ctx = s.setCtxRequestTimeout(ctx)
	ctx, parentSpan := s.StartRootSpan(ctx, constants.LOGOUT_ACCOUNT)
	defer parentSpan.Finish()

	if req == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	if req.Id == 0 {
		s.metrics.InvalidRequestParametersCounter.WithLabelValues(constants.LOGOUT_ACCOUNT).Inc()

		err := service_errors.ErrInvalidInputArguments
		var msg = "invalid input parameters. please specify a valid user id"
		s.logger.Error(err, msg)

		return nil, err
	}

	var (
		begin           = time.Now()
		took            = time.Since(begin)
		operation       = func() (interface{}, error) {
			return nil, s.authnClient.LogOutAccount()
		}
		retryableOperation = func() (interface{}, error) {
			return s.performRetryableRpcCall(ctx, operation)
		}
	)

	ctx = opentracing.ContextWithSpan(ctx, parentSpan)
	_, err := s.PerformRPCOperationAndInstrument(ctx, retryableOperation, constants.LOGOUT_ACCOUNT, &took)
	if err != nil {
		return nil, err
	}

	s.logger.For(ctx).Info("Successfully logged out user account", zap.Int("id", int(req.GetId())))
	response := &proto.LogoutAccountResponse{
		Error:                "",
	}
	return response, nil
}
