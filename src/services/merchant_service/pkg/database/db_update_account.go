package database

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"gorm.io/gorm"
)

// UpdateMerchantAccount updates a business account and saves it to the database as long as it exists
//
// the assumption from the context of the database is that all account should have the proper set of parameters in order prior
// to attempted storage. The client should handle any rpc operations to necessary prior to storage
func (db *Db) UpdateMerchantAccount(ctx context.Context, id uint64, account *merchant_service_proto_v1.MerchantAccount) (
	*merchant_service_proto_v1.MerchantAccount, error) {
	db.Logger.For(ctx).Info("creating business account")
	ctx, span := db.startRootSpan(ctx, "create_business_account_op")
	defer span.Finish()

	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		db.Logger.For(ctx).Info("starting transaction")
		span := db.TracingEngine.CreateChildSpan(ctx, "create_business_account_tx")
		defer span.Finish()

		if err := db.ValidateAccount(ctx, account); err != nil {
			return nil, err
		}

		if ok, err := db.FindMerchantAccountById(ctx, id); !ok && err != nil {
			return nil, errors.ErrAccountDoesNotExist
		}

		err := db.SaveAccountRecord(tx, account)
		if err != nil {
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
