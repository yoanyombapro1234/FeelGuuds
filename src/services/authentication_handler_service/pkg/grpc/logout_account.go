package grpc

import (
	"context"

	"go.uber.org/zap"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
)

// LogoutAccount revokes the user account session from the context of the authentication service
func (s *Server) LogoutAccount(ctx context.Context, req *proto.LogoutAccountRequest) (*proto.LogoutAccountResponse, error) {
	operationType := constants.LOGOUT_ACCOUNT
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

	_, err = s.PerformRetryableRPCOperation(ctx, rootSpan, callAuthenticationService, constants.LOGOUT_ACCOUNT)()
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	s.logger.Info("Successfully logged out user account", zap.Int("id", int(req.GetId())))
	response := &proto.LogoutAccountResponse{
		Error: "",
	}
	return response, nil
}
