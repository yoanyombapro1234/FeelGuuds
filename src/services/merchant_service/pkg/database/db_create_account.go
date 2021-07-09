package database

import (
	"context"
	"fmt"

	core_database "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-database"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"gorm.io/gorm"
)

// CreateMerchantAccount creates a business account and saves it to the database
// the assumption from the context of the database is that all account should have the proper set of parameters in order prior
// to attempted storage. The client should handle any rpc operations to necessary prior to storage
func (db *Db) CreateMerchantAccount(ctx context.Context, account *merchant_service_proto_v1.MerchantAccount) (*merchant_service_proto_v1.MerchantAccount, error) {
	const operation = "create_business_account_db_op"
	db.Logger.For(ctx).Info(fmt.Sprintf("create business account database operation."))
	ctx, span := db.startRootSpan(ctx, operation)
	defer span.Finish()

	tx := db.createAccountTxFunc(account)

	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	createdAccount := result.(*merchant_service_proto_v1.MerchantAccount)
	return createdAccount, nil
}

// createAccountTxFunc encloses the account creation step in a database transaction
func (db *Db) createAccountTxFunc(account *merchant_service_proto_v1.MerchantAccount) core_database.CmplxTx {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		const operation = "create_business_account_db_tx"
		db.Logger.For(ctx).Info("starting transaction")
		span := db.TracingEngine.CreateChildSpan(ctx, operation)
		defer span.Finish()

		if err := db.ValidateAccount(ctx, account); err != nil {
			return nil, err
		}

		// check if merchant account already exist
		if ok, err := db.FindMerchantAccountByEmail(ctx, account.BusinessEmail); ok && err == nil {
			return nil, errors.ErrAccountAlreadyExist
		}

		if err := db.UpdateAccountOnboardStatus(ctx, account); err != nil {
			return nil, err
		}

		if err := tx.Create(&account).Error; err != nil {
			return nil, err
		}

		return &account, nil
	}
	return tx
}
