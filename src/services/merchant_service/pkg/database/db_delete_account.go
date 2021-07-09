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
func (db *Db) DeleteMerchantAccount(ctx context.Context, id uint64) (bool, error) {
	db.Logger.For(ctx).Info("creating business account")
	ctx, span := db.startRootSpan(ctx, "delete_business_account_op")
	defer span.Finish()

	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		db.Logger.For(ctx).Info("starting transaction")
		span := db.TracingEngine.CreateChildSpan(ctx, "delete_business_account_tx")
		defer span.Finish()

		if id == 0 {
			return false, errors.ErrInvalidInputArguments
		}

		if ok, err := db.FindMerchantAccountById(ctx, id); !ok && err != nil {
			return false, err
		}

		tx = tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Model(&merchant_service_proto_v1.MerchantAccount{}).
			Where("id = ?", id)

		if err := tx.Update("account_state", merchant_service_proto_v1.MerchantAccountState_Inactive).Error; err != nil {
			db.Logger.For(ctx).Error(errors.ErrFailedToUpdateAccountActiveStatus, err.Error())
			return false, err
		}

		return true, nil
	}

	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return false, err
	}

	status, ok := result.(*bool)
	if !ok {
		return false, errors.ErrFailedToCastToType
	}

	return *status, nil
}
