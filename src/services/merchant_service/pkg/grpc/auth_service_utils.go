package grpc

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/authentication_handler_service/gen/proto"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
)

func (s *Server) CallAuthHandlerSvcAndAuthenticateAccount(ctx context.Context, merchantAcct *merchant_service_proto_v1.MerchantAccount) (*proto.
	AuthenticateAccountResponse, error) {
	authRpcReq := &proto.AuthenticateAccountRequest{Email: merchantAcct.BusinessEmail, Password: merchantAcct.Password}
	response, err := s.AuthenticationHandlerClient.AuthenticateAccount(ctx, authRpcReq)
	if err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return nil, err
	}
	return response, nil
}

func (s *Server) CallAuthHandlerSvcAndCreateAccount(ctx context.Context, merchantAcct *merchant_service_proto_v1.MerchantAccount) (
	*proto.CreateAccountResponse, error) {
	rpcReq := &proto.CreateAccountRequest{
		Email:    merchantAcct.BusinessEmail,
		Password: merchantAcct.Password,
	}
	authnAcct, err := s.AuthenticationHandlerClient.CreateAccount(ctx, rpcReq)
	if err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return nil, err
	}

	if authnAcct.Error != constants.EMPTY {
		rpcErr := errors.NewError(authnAcct.Error)
		s.logger.For(ctx).Error(rpcErr, rpcErr.Error())
		return nil, rpcErr
	}

	return authnAcct, err
}

func (s *Server) CallAuthHandlerSvcAndLockAccount(ctx context.Context, id uint32) error {
	rpcReq := &proto.LockAccountRequest{
		Id: id,
	}
	authnAcct, err := s.AuthenticationHandlerClient.LockAccount(ctx, rpcReq)
	if err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return err
	}

	if authnAcct.Error != constants.EMPTY {
		rpcErr := errors.NewError(authnAcct.Error)
		s.logger.For(ctx).Error(rpcErr, rpcErr.Error())
		return rpcErr
	}

	return err
}

func (s *Server) CallAuthHandlerSvcAndUnlockAccount(ctx context.Context, id uint32) error {
	rpcReq := &proto.UnLockAccountRequest{
		Id: id,
	}
	authnAcct, err := s.AuthenticationHandlerClient.UnLockAccount(ctx, rpcReq)
	if err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return err
	}

	if authnAcct.Error != constants.EMPTY {
		rpcErr := errors.NewError(authnAcct.Error)
		s.logger.For(ctx).Error(rpcErr, rpcErr.Error())
		return rpcErr
	}

	return err
}

func (s *Server) CallAuthHandlerSvcAndUpdateAccount(ctx context.Context, authnId uint64, email string) error {
	rpcReq := &proto.UpdateAccountRequest{
		Id:    uint32(authnId),
		Email: email,
	}

	authnAcct, err := s.AuthenticationHandlerClient.UpdateAccount(ctx, rpcReq)
	if err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return err
	}

	if authnAcct.Error != constants.EMPTY {
		rpcErr := errors.NewError(authnAcct.Error)
		s.logger.For(ctx).Error(rpcErr, rpcErr.Error())
		return rpcErr
	}

	return err
}
