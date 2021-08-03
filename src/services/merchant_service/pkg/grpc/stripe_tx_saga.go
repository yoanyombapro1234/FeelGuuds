package grpc

import (
	"context"

	"github.com/itimofeev/go-saga"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/service_errors"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/stripe_client"
)

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
		if resp.StripeId == constants.EMPTY{
			return service_errors.ErrInvalidInputArguments
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
