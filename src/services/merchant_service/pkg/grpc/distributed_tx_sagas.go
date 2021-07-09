package grpc

import (
	"context"

	"github.com/itimofeev/go-saga"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/stripe_client"
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

// sagaCreateAccountThroughStripe returns a saga comprised of an action and compensating function used to create an account through stripe or
// revert the created account in the face of failures
//
// The purpose of the action is to create an account in stripe while the purpose of the compensating function
// is to roll back any change performed in stripe if downstream saga actions fail
func (s *Server) sagaCreateAccountThroughStripe(request *merchant_service_proto_v1.CreateAccountRequest, merchantAcct *merchant_service_proto_v1.MerchantAccount, stripeResponseObj chan *stripe_client.Response) *saga.Step {
	return &saga.Step{
		Name: "create a connected account from the context of stripe",
		Func: func(stripeId chan<- *stripe_client.Response) func(ctx context.Context) error {
			return s.actionCreateAccountViaStripe(stripeId, request, merchantAcct)
		}(stripeResponseObj),
		CompensateFunc: func(stripeResponse <-chan *stripe_client.Response) func(ctx context.Context) error {
			return s.compensateDeleteAccountViaStripe(stripeResponse)
		}(stripeResponseObj),
		Options: nil,
	}
}

func (s *Server) compensateDeleteAccountViaStripe(stripeResponse <-chan *stripe_client.Response) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		resp := <-stripeResponse
		if resp.StripeId == 0 {
			return errors.ErrInvalidInputArguments
		}

		err := s.StripeClient.DeleteAccount(ctx, resp.StripeId)
		if err != nil {
			return err
		}
		return nil
	}
}

func (s *Server) actionCreateAccountViaStripe(stripeId chan<- *stripe_client.Response, request *merchant_service_proto_v1.CreateAccountRequest,
	merchantAcct *merchant_service_proto_v1.MerchantAccount) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// invoke and create connected account from context of stripe
		stripeResponse, err := s.StripeClient.CreateAccount(ctx, request)
		if err != nil {
			return err
		}

		merchantAcct.StripeConnectedAccountId = stripeResponse.StripeId
		stripeId <- stripeResponse
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
			return s.compensateDeleteAccountViaAuthHandlerSvc(authnAcctId)
		}(authnAcctId),
		Options: nil,
	}
}

func (s *Server) compensateDeleteAccountViaAuthHandlerSvc(authnAcctId chan uint32) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ctx, cancel := s.ConfigureOutgoingRpcRequest(ctx)
		defer cancel()

		id := <-authnAcctId
		if id == 0 {
			return errors.ErrInvalidInputArguments
		}

		err := s.CallAuthHandlerSvcAndDeleteAccount(ctx, id)
		if err != nil {
			return err
		}

		return nil
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
