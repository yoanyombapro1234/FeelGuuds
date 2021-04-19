package grpc

import (
	"context"

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
		s.logger.Error(err, "invalid input parameters. please specify a valid user id")

		return nil, err
	}

	var (
		operation = func() (interface{}, error) {
			return nil, s.authnClient.LogOutAccount()
		}
	)

	ctx = opentracing.ContextWithSpan(ctx, parentSpan)
	_, err := s.PerformRetryableRPCOperation(ctx, parentSpan, operation, constants.LOGOUT_ACCOUNT)()
	if err != nil {
		s.logger.Error(err, err.Error())
		return nil, err
	}

	s.logger.For(ctx).Info("Successfully logged out user account", zap.Int("id", int(req.GetId())))
	response := &proto.LogoutAccountResponse{
		Error: "",
	}
	return response, nil
}
