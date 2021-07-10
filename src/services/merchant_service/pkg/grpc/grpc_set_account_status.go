package grpc

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
)

func (s *Server) SetAccountStatus(ctx context.Context, request *merchant_service_proto_v1.SetAccountStatusRequest) (*merchant_service_proto_v1.SetAccountStatusResponse, error) {
	operationType := constants.UPDATE_MERCHANT_ACCOUNT_ONBOARDING_STATE
	ctx, rootSpan := s.ConfigureAndStartRootSpan(ctx, operationType)
	defer rootSpan.Finish()

	if request == nil || request.AccountId == 0 {
		s.logger.For(ctx).Error(errors.ErrInvalidInputArguments, errors.ErrInvalidInputArguments.Error())
		return nil, errors.ErrInvalidInputArguments
	}

	acct, err := s.DbConn.UpdateAccountOnboardingStatus(ctx, request.GetAccountId(), request.AccountState)
	if err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return nil, err
	}

	return &merchant_service_proto_v1.SetAccountStatusResponse{
		Account: acct,
	}, nil
}
