package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

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
		var msg = "invalid input parameters. please specify a username and password"
		s.logger.Error(err, msg)

		return nil, err
	}

	var (
		begin           = time.Now()
		took            = time.Since(begin)
		isAccountLocked = false
		operation       = func() (interface{}, error) {
			return s.authnClient.ImportAccount(req.Email, req.Password, isAccountLocked)
		}
		retryableOperation = func() (interface{}, error) {
			return s.performRetryableRpcCall(ctx, operation)
		}
	)

	ctx = opentracing.ContextWithSpan(ctx, parentSpan)
	result, err := s.PerformRPCOperationAndInstrument(ctx, retryableOperation, constants.CREATE_ACCOUNT, &took)
	if err != nil {
		return nil, err
	}

	authnID, ok := result.(int)
	if !ok {
		s.metrics.CastingOperationFailureCounter.WithLabelValues(constants.CREATE_ACCOUNT)
		err := errors.New("failed to convert result to uint32 id value")
		s.logger.For(ctx).ErrorM(err, "casting error")
		return nil, err
	}

	s.logger.For(ctx).Info("Successfully created user account", zap.Int("Id", int(authnID)))
	response := &proto.CreateAccountResponse{Id: uint32(authnID), Error: ""}

	return response, nil
}
