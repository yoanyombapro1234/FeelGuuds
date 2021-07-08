package grpc

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/pkg/service_errors"
	"go.uber.org/zap"
)

// AuthenticateAccount authenticates the current user account against the authentication service ensuring the credentials defined exist
func (s *Server) AuthenticateAccount(ctx context.Context, req *proto.AuthenticateAccountRequest) (*proto.AuthenticateAccountResponse, error) {
	operationType := constants.LOGIN_ACCOUNT
	ctx, rootSpan := s.ConfigureAndStartRootSpan(ctx, operationType)
	defer rootSpan.Finish()

	if req == nil {
		return nil, service_errors.ErrInvalidInputArguments
	}

	err, pwdOrEmailIsInvalid := s.IsPasswordOrEmailInValid(req.Email, req.Password, operationType)
	if pwdOrEmailIsInvalid {
		return nil, err
	}

	var callAuthenticationService = req.CallAuthenticationService(s.authnClient)

	result, err := s.PerformRetryableRPCOperation(ctx, rootSpan, callAuthenticationService, operationType)()
	if err != nil {
		s.logger.Error(err, err.Error())
		return nil, err
	}

	err, tokenIsInvalid, token := s.CheckJwtTokenForInValidity(ctx, result, operationType)
	if tokenIsInvalid {
		return nil, err
	}

	s.logger.For(ctx).Info("Successfully authenticated user account", zap.String("jwt", token))
	response := &proto.AuthenticateAccountResponse{
		Token: token,
		Error: "",
	}

	return response, nil
}
