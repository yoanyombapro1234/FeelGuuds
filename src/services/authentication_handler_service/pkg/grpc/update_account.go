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

// UpdateAccount updates the account via the authentication services
func (s *Server) UpdateAccount(ctx context.Context, req *proto.UpdateAccountRequest) (*proto.UpdateAccountResponse, error) {
	ctx = s.setCtxRequestTimeout(ctx)
	ctx, parentSpan := s.StartRootSpan(ctx, constants.UPDATE_ACCOUNT)
	defer parentSpan.Finish()

	if req == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	// TODO: change this to username instead of email
	if req.Email == "" || req.Id == 0 {
		s.metrics.InvalidRequestParametersCounter.WithLabelValues(constants.UPDATE_ACCOUNT).Inc()

		err := service_errors.ErrInvalidInputArguments
		var msg = "invalid input parameters. please specify an email or valid id"
		s.logger.Error(err, msg)

		return nil, err
	}

	var (
		begin           = time.Now()
		took            = time.Since(begin)
		operation       = func() (interface{}, error) {
			return nil, s.authnClient.Update(strconv.Itoa(int(req.Id)), req.Email)
		}
		retryableOperation = func() (interface{}, error) {
			return s.performRetryableRpcCall(ctx, operation)
		}
	)

	ctx = opentracing.ContextWithSpan(ctx, parentSpan)
	_, err := s.PerformRPCOperationAndInstrument(ctx, retryableOperation, constants.UPDATE_ACCOUNT, &took)
	if err != nil {
		return nil, err
	}

	s.logger.For(ctx).Info("Successfully updated user account", zap.Int("Id", int(req.Id)))
	response := &proto.UpdateAccountResponse{ Error: ""}
	return response, nil
}
