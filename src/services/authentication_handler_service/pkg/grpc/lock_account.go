package grpc

import (
	"context"

	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

// LockAccount locks an account as long as it exists from the context of the authentication service
func (s *Server) LockAccount(ctx context.Context, req *proto.LockAccountRequest) (*proto.LockAccountResponse, error) {
	operationType := constants.LOCK_ACCOUNT
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

	_, err = s.PerformRetryableRPCOperation(ctx, rootSpan, callAuthenticationService, constants.LOCK_ACCOUNT)()
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	s.logger.Info("Successfully locked user account", zap.Int("Id", int(req.GetId())))
	return &proto.LockAccountResponse{
		Error: "",
	}, nil
}
