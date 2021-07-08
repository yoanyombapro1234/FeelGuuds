package database

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"gorm.io/gorm"
)

// GetMerchantAccountsById obtains a set of business accounts by specified ids
//
// the assumption from the context of the database is that all account should have the proper set of parameters in order prior
// to attempted storage. The client should handle any rpc operations to necessary prior to storage
func (db *Db) GetMerchantAccountsById(ctx context.Context, ids []uint64) ([]*merchant_service_proto_v1.MerchantAccount, error) {
	db.Logger.For(ctx).Info("creating business account")
	ctx, span := db.startRootSpan(ctx, "create_business_account_op")
	defer span.Finish()

	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		db.Logger.For(ctx).Info("starting transaction")
		span := db.TracingEngine.CreateChildSpan(ctx, "create_business_account_tx")
		defer span.Finish()

		var accounts = make([]*merchant_service_proto_v1.MerchantAccount, len(ids)+1)
		if err := tx.Where(ids).Find(&accounts).Error; err != nil {
			return nil, err
		}

		return accounts, nil
	}

	result, err := db.Conn.PerformComplexTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	accounts := result.([]*merchant_service_proto_v1.MerchantAccount)
	return accounts, nil
}
