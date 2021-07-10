package database

import (
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"gorm.io/gorm"
)

// SaveAccountRecord saves a record in the database
func (db *Db) SaveAccountRecord(tx *gorm.DB, account *merchant_service_proto_v1.MerchantAccount) error {
	pswd, err := db.ValidateAndHashPassword(account.Password)
	if err != nil {
		return err
	}

	account.Password = pswd
	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&account).Error; err != nil {
		return err
	}
	return nil
}
