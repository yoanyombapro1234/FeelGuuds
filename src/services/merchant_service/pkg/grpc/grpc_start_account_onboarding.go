package grpc

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
)

func (s *Server) StartAccountOnboarding(ctx context.Context, request *merchant_service_proto_v1.StartAccountOnboardingRequest) (*merchant_service_proto_v1.StartAccountOnboardingRespone, error) {
	// this endpoint will be invoked if the refresh url sent to stripe is hit by the client
	// perform request validations
	operationType := constants.START_MERCHANT_ACCOUNT_ONBOARDING
	ctx, rootSpan := s.ConfigureAndStartRootSpan(ctx, operationType)
	defer rootSpan.Finish()

	// check if the account exists in the database
	acc, err := s.DbConn.GetMerchantAccountById(ctx, request.GetAccountId())
	if err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return nil, err
	}

	// call the stripe api for the url to which the account will route to
	response, err := s.StripeClient.GetAccountLink(ctx, acc.StripeConnectedAccountId)
	if err != nil {
		return nil, err
	}

	return &merchant_service_proto_v1.StartAccountOnboardingRespone{Url: response.Url}, nil
}
