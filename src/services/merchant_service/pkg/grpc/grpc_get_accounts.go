package grpc

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
)

func (s *Server) GetAccounts(ctx context.Context, request *merchant_service_proto_v1.GetAccountsRequest) (*merchant_service_proto_v1.GetAccountsResponse, error) {
	operationType := constants.GET_MERCHANT_ACCOUNTS
	ctx, rootSpan := s.ConfigureAndStartRootSpan(ctx, operationType)
	defer rootSpan.Finish()

	if request == nil || len(request.AccountIds) == 0 {
		s.logger.For(ctx).Error(errors.ErrInvalidInputArguments, errors.ErrInvalidInputArguments.Error())
		return nil, errors.ErrInvalidInputArguments
	}

	accts, err := s.DbConn.GetMerchantAccountsById(ctx, request.AccountIds)
	if err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return nil, err
	}

	return &merchant_service_proto_v1.GetAccountsResponse{Accounts: accts}, nil
}
