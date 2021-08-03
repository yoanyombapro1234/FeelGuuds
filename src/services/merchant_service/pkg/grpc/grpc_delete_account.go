package grpc

import (
	"context"

	"github.com/itimofeev/go-saga"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/service_errors"
)

func (s *Server) DeleteAccount(ctx context.Context, request *merchant_service_proto_v1.DeleteAccountRequest) (*merchant_service_proto_v1.DeleteAccountResponse, error) {
	// perform request validations
	operationType := constants.DELETE_MERCHANT_ACCOUNT
	ctx, rootSpan := s.ConfigureAndStartRootSpan(ctx, operationType)
	defer rootSpan.Finish()

	if request == nil || request.AccountId == 0 {
		s.logger.For(ctx).Error(service_errors.ErrInvalidInputArguments, service_errors.ErrInvalidInputArguments.Error())
		return nil, service_errors.ErrInvalidInputArguments
	}

	/*
		the feelguuds platform only performs soft deletes. Hence to do so, we invoke only the authentication service.
		we perform these operations as distributed tx by use of sagas.
		sage logic flow:
		1. lock the account from the context of the authentication service
		2. deactivate the account from the context of the local database
	*/

	// check if account exists in local db
	acct, err := s.DbConn.GetMerchantAccountById(ctx, request.GetAccountId())
	if err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return nil, err
	}

	//  dispatch set of sagas
	var authnAcctId = make(chan uint32, 1)
	var merchantAccId = make(chan uint64, 1)
	sagaSteps := make([]*saga.Step, 0)
	merchantAccId <- acct.GetId()

	lockAccountStep := s.sagaLockAccountThroughAuthenticationHandlerService(authnAcctId)
	deleteRecordFromDbStep := s.sagaDeactivateAccountInDB(merchantAccId)
	sagaSteps = append(sagaSteps, lockAccountStep, deleteRecordFromDbStep)
	if err := s.DbConn.Saga.RunSaga(ctx, "delete_account", sagaSteps...); err != nil {
		s.logger.For(ctx).Error(err, err.Error())
		return nil, err
	}

	return &merchant_service_proto_v1.DeleteAccountResponse{
		IsDeleted: true,
	}, nil
}
