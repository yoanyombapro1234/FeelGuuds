package grpc

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
	"go.uber.org/zap"
)

// CreateAccount creates a user account via the authentication service
func (s *Server) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	operationType := constants.CREATE_ACCOUNT
	ctx, rootSpan := s.ConfigureAndStartRootSpan(ctx, operationType)
	defer rootSpan.Finish()

	if req == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	email, password := req.Email, req.Password
	err, pwdOrEmailIsInvalid := s.IsPasswordOrEmailInValid(email, password, operationType)
	if pwdOrEmailIsInvalid {
		return nil, err
	}

	var callAuthenticationService = req.CallAuthenticationService(s.authnClient)

	result, err := s.PerformRetryableRPCOperation(ctx, rootSpan, callAuthenticationService, operationType)()
	if err != nil {
		return nil, err
	}

	id, err := s.GetIdFromResponseObject(ctx, result, operationType)
	if err != nil {
		return nil, err
	}

	s.logger.For(ctx).Info("Successfully created user account", zap.Int("Id", int(id)))
	response := &proto.CreateAccountResponse{Id: uint32(id), Error: ""}

	return response, nil
}

