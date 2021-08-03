package grpc

import (
	"context"

	"github.com/itimofeev/go-saga"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/service_errors"
)

// sagaGetJwtTokenFromAuthHandlerSvc returns a saga comprised of an action and compensating function used to authenticate an account through the
// authentication handler service
//
// The purpose of the action is to authenticate an account in the auth handler service while the purpose of the compensating function
// is to roll back any change performed if downstream saga actions fail
func (s *Server) sagaGetJwtTokenFromAuthHandlerSvc(merchantAcct *merchant_service_proto_v1.MerchantAccount, jwtToken chan string) *saga.Step {
	return &saga.Step{
		Name: "get jwt token from authentication handler service",
		Func: func(jwtToken chan<- string) func(ctx context.Context) error {
			return s.actionGetJwtTokenViaAuthHandlerSvc(jwtToken, merchantAcct)
		}(jwtToken),
		CompensateFunc: func(ctx context.Context) error {
			return nil
		},
		Options: nil,
	}
}

func (s *Server) actionGetJwtTokenViaAuthHandlerSvc(jwtToken chan<- string, merchantAcct *merchant_service_proto_v1.MerchantAccount) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ctx, cancel := s.ConfigureOutgoingRpcRequest(ctx)
		defer cancel()

		response, err := s.CallAuthHandlerSvcAndAuthenticateAccount(ctx, merchantAcct)
		if err != nil {
			return err
		}

		jwtToken <- response.Token
		return nil
	}
}

// sagaCreateAccountThroughAuthenticationHandlerService returns a saga comprised of an action and compensating function.
//
// The purpose of the action is to create an account in the authentication handler service while the purpose of the compensating function
// is to roll back any change performed in the authentication handler service if downstream saga actions fail
func (s *Server) sagaCreateAccountThroughAuthenticationHandlerService(merchantAcct *merchant_service_proto_v1.MerchantAccount, authnAcctId chan uint32) *saga.Step {
	return &saga.Step{
		Name: "create account in authentication handler service",
		Func: func(authnId chan<- uint32) func(ctx context.Context) error {
			return s.actionCreateAccountViaAuthHandlerSvc(authnId, merchantAcct)
		}(authnAcctId),
		CompensateFunc: func(authnId <-chan uint32) func(ctx context.Context) error {
			// delete the account from the context of the authentication service if errors arise
			return s.compensateLockAccountViaAuthHandlerSvc(authnAcctId)
		}(authnAcctId),
		Options: nil,
	}
}

func (s *Server) compensateLockAccountViaAuthHandlerSvc(authnAcctId chan uint32) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ctx, cancel := s.ConfigureOutgoingRpcRequest(ctx)
		defer cancel()

		id := <-authnAcctId
		if id == 0 {
			return service_errors.ErrInvalidInputArguments
		}

		return s.CallAuthHandlerSvcAndLockAccount(ctx, id)
	}
}

func (s *Server) actionCreateAccountViaAuthHandlerSvc(authnId chan<- uint32, merchantAcct *merchant_service_proto_v1.MerchantAccount) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ctx, cancel := s.ConfigureOutgoingRpcRequest(ctx)
		defer cancel()

		authnAcct, err := s.CallAuthHandlerSvcAndCreateAccount(ctx, merchantAcct)
		if err != nil {
			return err
		}

		// update the merchant account - authn id with the rpc response
		merchantAcct.AuthnAccountId = uint64(authnAcct.Id)
		authnId <- authnAcct.Id
		return nil
	}
}

// sagaLockAccountThroughAuthenticationHandlerService returns a saga comprised of an action and compensating function.
//
// The purpose of the action is to lock an account in the authentication handler service while the purpose of the compensating function
// is to roll back any change performed in the authentication handler service if downstream saga actions fail
func (s *Server) sagaLockAccountThroughAuthenticationHandlerService(
	authnAcctId chan uint32) *saga.Step {
	return &saga.Step{
		Name: "delete account in authentication handler service",
		Func: func(authnId <-chan uint32) func(ctx context.Context) error {
			// delete the account from the context of the authentication service if errors arise
			return s.compensateLockAccountViaAuthHandlerSvc(authnAcctId)
		}(authnAcctId),
		CompensateFunc: func(authnId <-chan uint32) func(ctx context.Context) error {
			// unlock the account from the context of the authentication service if errors arise
			return s.compensateUnLockAccountViaAuthHandlerSvc(authnAcctId)
		}(authnAcctId),
		Options: nil,
	}
}

func (s *Server) compensateUnLockAccountViaAuthHandlerSvc(authnAcctId chan uint32) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ctx, cancel := s.ConfigureOutgoingRpcRequest(ctx)
		defer cancel()

		id := <-authnAcctId
		if id == 0 {
			return service_errors.ErrInvalidInputArguments
		}

		return s.CallAuthHandlerSvcAndUnlockAccount(ctx, id)
	}
}

// sagaUpdateAccountThroughAuthenticationHandlerService returns a saga comprised of an action and compensating function.
//
// The purpose of the action is to update an account in the authentication handler service while the purpose of the compensating function
// is to roll back any change performed in the authentication handler service if downstream saga actions fail
func (s *Server) sagaUpdateAccountThroughAuthenticationHandlerService(
	authnAcctId uint64, oldEmail, newEmail string) *saga.Step {
	return &saga.Step{
		Name: "update account email in authentication handler service",
		Func: func(authnId uint64, email string) func(ctx context.Context) error {
			return s.updateAccountEmailViaAuthHandlerSvc(authnId, email)
		}(authnAcctId, newEmail),
		CompensateFunc: func(authnId uint64, email string) func(ctx context.Context) error {
			return s.updateAccountEmailViaAuthHandlerSvc(authnId, email)
		}(authnAcctId, oldEmail),
		Options: nil,
	}
}

func (s *Server) updateAccountEmailViaAuthHandlerSvc(id uint64, email string) func(
	ctx context.Context) error {
	return func(ctx context.Context) error {
		ctx, cancel := s.ConfigureOutgoingRpcRequest(ctx)
		defer cancel()

		return s.CallAuthHandlerSvcAndUpdateAccount(ctx, id, email)
	}
}
