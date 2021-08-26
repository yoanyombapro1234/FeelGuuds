package grpc

import (
	"context"

	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

// UnLockAccount unlocks a user account from the context of the authentication service as long as the account exists
func (s *Server) UnLockAccount(ctx context.Context, req *proto.UnLockAccountRequest) (*proto.UnLockAccountResponse, error) {
	operationType := constants.UNLOCK_ACCOUNT
	ctx, rootSpan := s.ConfigureAndStartRootSpan(ctx, operationType)
	defer rootSpan.Finish()

	if req == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	err, ok := s.IsValidID(req.Id, operationType)
	if !ok {
		return nil, err
	}

	var callAuthenticationService = req.CallAuthenticationService(s.authnClient)

	_, err = s.PerformRetryableRPCOperation(ctx, rootSpan, callAuthenticationService, constants.UNLOCK_ACCOUNT)()
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	s.logger.Info("Successfully unlocked user account", zap.Int("Id", int(req.GetId())))
	return &proto.UnLockAccountResponse{
		Error: "",
	}, nil
}
