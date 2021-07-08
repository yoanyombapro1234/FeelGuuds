package database

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetMerchantAccountById finds a merchant account by id
func (db *Db) GetMerchantAccountById(ctx context.Context, id uint64) (*merchant_service_proto_v1.MerchantAccount, error) {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		span := db.TracingEngine.CreateChildSpan(ctx, "get_merchant_account_by_id_tx")
		defer span.Finish()

		if id == 0 {
			return false, errors.ErrInvalidInputArguments
		}

		var account merchant_service_proto_v1.MerchantAccount
		if err := tx.Preload(clause.Associations).Where(&merchant_service_proto_v1.MerchantAccount{Id: id}).First(&account).Error; err != nil {
			return true, errors.ErrAccountDoesNotExist
		}

		return false, nil
	}

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
