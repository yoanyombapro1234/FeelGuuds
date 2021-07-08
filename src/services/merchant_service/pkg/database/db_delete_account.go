package database

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DeleteMerchantAccount deletes a business account and updates the database
//
// the assumption from the context of the database is that all account should have the proper set of parameters in order prior
// to attempted storage. The client should handle any rpc operations to necessary prior to storage
func (db *Db) DeleteMerchantAccount(ctx context.Context, account *merchant_service_proto_v1.MerchantAccount) (*merchant_service_proto_v1.
	MerchantAccount, error) {
	db.Logger.For(ctx).Info("creating business account")
	ctx, span := db.startRootSpan(ctx, "delete_business_account_op")
	defer span.Finish()

	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		db.Logger.For(ctx).Info("starting transaction")
		span := db.TracingEngine.CreateChildSpan(ctx, "delete_business_account_tx")
		defer span.Finish()

		if err := db.ValidateAccount(ctx, account); err != nil {
			return nil, err
		}

		if ok, err := db.FindMerchantAccountById(ctx, account.Id); !ok && err != nil {
			return nil, err
		}

		tx = tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Model(&merchant_service_proto_v1.MerchantAccount{}).
			Where("id = ?", account.Id)

		if err := tx.Update("account_state", merchant_service_proto_v1.MerchantAccountState_Inactive).Error; err != nil {
			db.Logger.For(ctx).Error(errors.ErrFailedToUpdateAccountActiveStatus, err.Error())
			return nil, err
		}

		return &account, nil
	}

	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	createdAccount := result.(*merchant_service_proto_v1.MerchantAccount)
	return createdAccount, nil
}
