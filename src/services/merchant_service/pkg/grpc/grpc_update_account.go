package grpc

import (
	"context"

	"github.com/itimofeev/go-saga"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
)

func (s *Server) UpdateAccount(ctx context.Context, request *merchant_service_proto_v1.UpdateAccountRequest) (*merchant_service_proto_v1.UpdateAccountResponse, error) {
	sagaSteps := make([]*saga.Step, 0)
	var merchantAcctId = make(chan uint64, 1)

	// perform request validations
	operationType := constants.UPDATE_MERCHANT_ACCOUNT
	ctx, rootSpan := s.ConfigureAndStartRootSpan(ctx, operationType)
	defer rootSpan.Finish()

	if request == nil || request.Account == nil {
		s.logger.For(ctx).Error(errors.ErrInvalidInputArguments, errors.ErrInvalidInputArguments.Error())
		return nil, errors.ErrInvalidInputArguments
	}

	// update scenarios ....
	// 1. update all fields except the email and password
	// 2. update the email or password field
	acct, err := s.DbConn.GetMerchantAccountById(ctx, request.GetAccountId())
	if err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return nil, err
	}

	if PasswordAndEmailUnchanged(acct, request) {
		// save the record in the database
		acc, err := s.DbConn.UpdateMerchantAccount(ctx, request.AccountId, request.Account)
		if err != nil {
			s.logger.For(ctx).Error(err, err.Error())
			return nil, err
		}

		return &merchant_service_proto_v1.UpdateAccountResponse{Account: acc}, nil
	}

	// populate the channel
	merchantAcctId <- acct.Id

	// perform the distributed update tx as a set of saga
	updateEmailViaAuthSvc := s.sagaUpdateAccountThroughAuthenticationHandlerService(acct.AuthnAccountId, acct.BusinessEmail,
		request.Account.BusinessEmail)
	updateEmailViaDb := s.sagaSaveCreatedAccountInDB(acct, merchantAcctId)

	sagaSteps = append(sagaSteps, updateEmailViaAuthSvc, updateEmailViaDb)
	if err := s.DbConn.Saga.RunSaga(ctx, "update_account", sagaSteps...); err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return nil, err
	}

	return &merchant_service_proto_v1.UpdateAccountResponse{Account: acct}, nil
}

func PasswordAndEmailUnchanged(acct *merchant_service_proto_v1.MerchantAccount, request *merchant_service_proto_v1.UpdateAccountRequest) bool {
	return acct.BusinessEmail == request.Account.BusinessEmail
}
