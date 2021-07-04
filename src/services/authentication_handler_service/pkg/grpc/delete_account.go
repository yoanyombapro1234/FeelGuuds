package grpc

import (
	"context"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

// DeleteAccount deletes a user account via the authentication service
func (s *Server) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.DeleteAccountResponse, error) {
	ctx = s.setCtxRequestTimeout(ctx)
	ctx, parentSpan := s.StartRootSpan(ctx, constants.DELETE_ACCOUNT)
	defer parentSpan.Finish()

	if req == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	if req.Id == 0 {
		s.metrics.InvalidRequestParametersCounter.WithLabelValues(constants.DELETE_ACCOUNT).Inc()
		err := service_errors.ErrInvalidInputArguments
		s.logger.Error(err, "invalid input parameters. please specify a valid user id")
		return nil, err
	}

	var (
		operation = func() (interface{}, error) {
			if err := s.authnClient.ArchiveAccount(strconv.Itoa(int(req.GetId()))); err != nil {
				return nil, err
			}
			return nil, nil
		}
	)

	ctx = opentracing.ContextWithSpan(ctx, parentSpan)
	_, err := s.PerformRetryableRPCOperation(ctx, parentSpan, operation, constants.DELETE_ACCOUNT)()
	if err != nil {
		s.logger.Error(err, err.Error())
		return nil, err
	}

	s.logger.For(ctx).Info("Successfully archived user account", zap.Int("Id", int(req.GetId())))
	return &proto.DeleteAccountResponse{
		Error: "",
	}, nil
}
