package grpc

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/service_errors"
)

func (s *Server) GetAccount(ctx context.Context, request *merchant_service_proto_v1.GetAccountRequest) (*merchant_service_proto_v1.GetAccountResponse, error) {
	operationType := constants.GET_MERCHANT_ACCOUNT
	ctx, rootSpan := s.ConfigureAndStartRootSpan(ctx, operationType)
	defer rootSpan.Finish()

	if request == nil || request.AccountId == 0 {
		s.logger.For(ctx).Error(service_errors.ErrInvalidInputArguments, service_errors.ErrInvalidInputArguments.Error())
		return nil, service_errors.ErrInvalidInputArguments
	}

	acct, err := s.DbConn.GetMerchantAccountById(ctx, request.GetAccountId())
	if err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return nil, err
	}

	return &merchant_service_proto_v1.GetAccountResponse{
		Account: acct,
	}, nil
}
