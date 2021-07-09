package database

import (
	"context"
	"fmt"

	core_database "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-database"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetMerchantAccountById finds a merchant account by id
func (db *Db) GetMerchantAccountById(ctx context.Context, id uint64) (*merchant_service_proto_v1.MerchantAccount, error) {
	const operation = "get_business_account_by_id_db_op"
	db.Logger.For(ctx).Info(fmt.Sprintf("get business account by id database operation. id : %d", id))
	ctx, span := db.startRootSpan(ctx, operation)
	defer span.Finish()

	tx := db.getMerchantAccountByIdTxFunc(id)
	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	acc, ok := result.(*merchant_service_proto_v1.MerchantAccount)
	if !ok {
		return nil, errors.ErrFailedToCastToType
	}

	return acc, nil
}

// getMerchantAccountByIdTxFunc gets the merchant account by id operation wrapped in a database transaction.
func (db *Db) getMerchantAccountByIdTxFunc(id uint64) core_database.CmplxTx {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operation = "get_business_account_by_id_db_tx"
		span := db.TracingEngine.CreateChildSpan(ctx, operation)
		defer span.Finish()

		if id == 0 {
			return nil, errors.ErrInvalidInputArguments
		}

		var account merchant_service_proto_v1.MerchantAccount
		if err := tx.Preload(clause.Associations).Where(&merchant_service_proto_v1.MerchantAccount{Id: id}).First(&account).Error; err != nil {
			return nil, errors.ErrAccountDoesNotExist
		}

		if account.GetAccountState() == merchant_service_proto_v1.MerchantAccountState_Inactive {
			return nil, errors.ErrAccountDoesNotExist
		}

		return &account, nil
	}
	return tx
}
