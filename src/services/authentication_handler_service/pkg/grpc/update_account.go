package grpc

import (
	"context"

	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

// UpdateAccount updates the account via the authentication services
func (s *Server) UpdateAccount(ctx context.Context, req *proto.UpdateAccountRequest) (*proto.UpdateAccountResponse, error) {
	operationType := constants.UPDATE_ACCOUNT
	ctx, rootSpan := s.ConfigureAndStartRootSpan(ctx, operationType)
	defer rootSpan.Finish()

	if req == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	err, ok := s.IsValidID(req.Id, operationType)
	if!ok {
		return nil, err
	}

	err, ok = s.IsValidEmail(req.Email, operationType)
	if !ok {
		return nil, err
	}

	var callAuthenticationService = req.CallAuthenticationService(s.authnClient)

	_, err = s.PerformRetryableRPCOperation(ctx, rootSpan, callAuthenticationService, constants.GET_ACCOUNT)()
	if err != nil {
		s.logger.Error(err, err.Error())
		return nil, err
	}

	s.logger.For(ctx).Info("Successfully updated user account", zap.Int("Id", int(req.Id)))
	response := &proto.UpdateAccountResponse{Error: ""}
	return response, nil
}
