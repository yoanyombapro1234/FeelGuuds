package grpc

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

// CreateAccount creates a user account via the authentication service
func (s *Server) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	ctx = s.setCtxRequestTimeout(ctx)
	ctx, parentSpan := s.StartRootSpan(ctx, constants.CREATE_ACCOUNT)
	defer parentSpan.Finish()

	if req == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	if req.Email == "" || req.Password == "" {
		s.metrics.InvalidRequestParametersCounter.WithLabelValues(constants.CREATE_ACCOUNT).Inc()

		err := service_errors.ErrInvalidInputArguments
		s.logger.Error(err, "invalid input parameters. please specify a valid username and password")

		return nil, err
	}

	var (
		isAccountLocked = false
		operation       = func() (interface{}, error) {
			return s.authnClient.ImportAccount(req.Email, req.Password, isAccountLocked)
		}
	)

	ctx = opentracing.ContextWithSpan(ctx, parentSpan)
	result, err := s.PerformRetryableRPCOperation(ctx, parentSpan, operation, constants.CREATE_ACCOUNT)()
	if err != nil {
		return nil, err
	}

	id, ok := result.(int)
	if !ok {
		s.metrics.CastingOperationFailureCounter.WithLabelValues(constants.CREATE_ACCOUNT)
		err := status.Errorf(codes.Internal, "failed to convert result to uint32 id value")
		s.logger.For(ctx).ErrorM(err, "casting error")
		return nil, err
	}

	s.logger.For(ctx).Info("Successfully created user account", zap.Int("Id", int(id)))
	response := &proto.CreateAccountResponse{Id: uint32(id), Error: ""}

	return response, nil
}
