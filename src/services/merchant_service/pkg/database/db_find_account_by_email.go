package database

import (
	"context"
	"fmt"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"gorm.io/gorm"
)

// FindMerchantAccountByEmail finds a merchant account by email
func (db *Db) FindMerchantAccountByEmail(ctx context.Context, email string) (bool, error) {
	const operation = "merchant_account_exists_by_email_db_op"
	db.Logger.For(ctx).Info(fmt.Sprintf("get business account by email database operation."))
	ctx, span := db.startRootSpan(ctx, operation)
	defer span.Finish()

	tx := db.findMerchantAccountByEmailTxFunc(email)
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

// findMerchantAccountByEmailTxFunc wraps the logic in a db tx and returns it
func (db *Db) findMerchantAccountByEmailTxFunc(email string) func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operation = "merchant_account_exists_by_email_tx"
		db.Logger.For(ctx).Info(fmt.Sprintf("get business account by email database tx."))
		ctx, span := db.startRootSpan(ctx, operation)
		defer span.Finish()

		if email == constants.EMPTY {
			return false, errors.ErrInvalidInputArguments
		}

		var account merchant_service_proto_v1.MerchantAccount
		if err := tx.Where(&merchant_service_proto_v1.MerchantAccount{BusinessEmail: email}).First(&account).Error; err != nil {
			return false, errors.ErrAccountDoesNotExist
		}

		if ok := db.AccountActive(&account); !ok {
			return false, errors.ErrAccountDoesNotExist
		}

		return true, nil
	}
	return tx
}
