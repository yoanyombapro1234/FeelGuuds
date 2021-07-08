package database

import (
	"context"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"gorm.io/gorm"
)

// DoesMerchantAccountExist checks if a merchant account exists solely off its Id
func (db *Db) DoesMerchantAccountExist(ctx context.Context, id uint64) (bool, error) {
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		span := db.TracingEngine.CreateChildSpan(ctx, "does_merchant_account_exists_by_id_tx")
		defer span.Finish()

		if id == 0 {
			return false, errors.ErrInvalidInputArguments
		}

		if ok, err := db.FindMerchantAccountById(ctx, id); !ok && err != nil {
			return false, errors.ErrAccountDoesNotExist
		}

		return true, nil
	}

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
