package grpc

import (
	"context"

	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

// GetAccount obtains an account as long as the account exists from the context of the authentication service
func (s *Server) GetAccount(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
	operationType := constants.GET_ACCOUNT
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

	result, err := s.PerformRetryableRPCOperation(ctx, rootSpan, callAuthenticationService, constants.GET_ACCOUNT)()
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	account, err := s.GetAccountFromResponseObject(ctx, ok, result, operationType)
	if err != nil {
		return nil, err
	}

	s.logger.Info("Successfully obtained user account", zap.Int("Id", int(req.GetId())))
	return &proto.GetAccountResponse{
		Account: &proto.Account{
			Id:            uint32(account.ID),
			Username:      account.Username,
			Locked:        account.Locked,
			Deleted:       account.Deleted,
			XXX_sizecache: 0,
		},
		Error: "",
	}, nil
}
