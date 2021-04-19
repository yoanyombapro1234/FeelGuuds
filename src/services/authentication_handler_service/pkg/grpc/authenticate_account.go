package grpc

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

// AuthenticateAccount authenticates the current user account against the authentication service ensuring the credentials defined exist
func (s *Server) AuthenticateAccount(ctx context.Context, req *proto.AuthenticateAccountRequest) (*proto.AuthenticateAccountResponse, error) {
	ctx = s.setCtxRequestTimeout(ctx)
	ctx, parentSpan := s.StartRootSpan(ctx, constants.LOGIN_ACCOUNT)
	defer parentSpan.Finish()

	if req == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	if req.Email == "" || req.Password == "" {
		s.metrics.InvalidRequestParametersCounter.WithLabelValues(constants.LOGIN_ACCOUNT).Inc()

		err := service_errors.ErrInvalidInputArguments
		s.logger.Error(err, "invalid input parameters. please specify a valid email and password")

		return nil, err
	}

	var (
		operation = func() (interface{}, error) {
			return s.authnClient.LoginAccount(req.Email, req.Password)
		}
	)

	ctx = opentracing.ContextWithSpan(ctx, parentSpan)
	result, err := s.PerformRetryableRPCOperation(ctx, parentSpan, operation, constants.LOGIN_ACCOUNT)()
	if err != nil {
		s.logger.Error(err, err.Error())
		return nil, err
	}

	token, ok := result.(string)
	if !ok {
		s.metrics.CastingOperationFailureCounter.WithLabelValues(constants.LOGIN_ACCOUNT)
		err := errors.New("issue casting to jwt token")
		s.logger.For(ctx).ErrorM(err, "casting error")
		return nil, err
	}

	s.logger.For(ctx).Info("Successfully authenticated user account", zap.String("jwt", token))
	response := &proto.AuthenticateAccountResponse{
		Token: token,
		Error: "",
	}

	return response, nil
}
