package grpc

import (
	"context"

	"github.com/itimofeev/go-saga"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
)

func (s *Server) sagaSaveCreatedAccountThroughDB(merchantAcct *merchant_service_proto_v1.MerchantAccount, merchantAcctId chan uint64) *saga.Step {
	saveAccountToDb := &saga.Step{
		Name: "save account record in database",
		Func: func(merchantAcctId chan<- uint64) func(ctx context.Context) error {
			return s.actionSaveCreatedAccountInDb(merchantAcctId, merchantAcct)
		}(merchantAcctId),
		CompensateFunc: func(merchantAcctId <-chan uint64) func(ctx context.Context) error {
			return s.compensateDeleteCreatedAccountInDb(merchantAcctId)
		}(merchantAcctId),
		Options: nil,
	}
	return saveAccountToDb
}

func (s *Server) compensateDeleteCreatedAccountInDb(merchantAcctId <-chan uint64) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		id := <-merchantAcctId
		ok, err := s.DbConn.DeleteMerchantAccount(ctx, uint64(id))
		if !ok && err != nil {
			s.logger.For(ctx).Error(err, err.Error())
			return err
		}

		return nil
	}
}

func (s *Server) actionSaveCreatedAccountInDb(merchantAcctId chan<- uint64, merchantAcct *merchant_service_proto_v1.MerchantAccount) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// 4. store record in database after applying mutations
		newAcct, err := s.DbConn.CreateMerchantAccount(ctx, merchantAcct)
		if err != nil {
			s.logger.For(ctx).Error(err, err.Error())
			return err
		}
		merchantAcctId <- newAcct.Id
		return nil
	}
}
