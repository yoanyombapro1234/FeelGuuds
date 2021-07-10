package database

import (
	"context"
	"fmt"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"gorm.io/gorm"
)

// FindMerchantAccountById finds a merchant account by id
func (db *Db) FindMerchantAccountById(ctx context.Context, id uint64) (bool, error) {
	const operation = "merchant_account_exists_by_id_op"
	db.Logger.For(ctx).Info(fmt.Sprintf("get business account by id database operation."))
	ctx, span := db.startRootSpan(ctx, operation)
	defer span.Finish()

	tx := db.findMerchantAccountByIdTxFunc(id)
	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return true, err
	}

	status, ok := result.(*bool)
	if !ok {
		return true, errors.ErrFailedToCastToType
	}

	return *status, nil
}

// findMerchantAccountByIdTxFunc finds the merchant account by id and wraps it in a db tx.
func (db *Db) findMerchantAccountByIdTxFunc(id uint64) func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
	return func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operation = "merchant_account_exists_by_id_tx"
		db.Logger.For(ctx).Info(fmt.Sprintf("get business account by id database tx."))
		ctx, span := db.startRootSpan(ctx, operation)
		defer span.Finish()

		if id == 0 {
			return false, errors.ErrInvalidInputArguments
		}

		var account merchant_service_proto_v1.MerchantAccount
		if err := tx.Where(&merchant_service_proto_v1.MerchantAccount{Id: id}).First(&account).Error; err != nil {
			return false, errors.ErrAccountDoesNotExist
		}

		if ok := db.AccountActive(&account); !ok {
			return false, errors.ErrAccountDoesNotExist
		}

		return true, nil
	}
}
