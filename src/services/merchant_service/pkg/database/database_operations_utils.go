package database

import (
	"context"
	"fmt"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"gorm.io/gorm"
)

// ValidateAccount performs various account level validations
func (db *Db) ValidateAccount(ctx context.Context, account *merchant_service_proto_v1.MerchantAccount) error {
	err := db.ValidateAccountNotNil(ctx, account)
	if err != nil {
		db.Logger.Error(err, err.Error())
		return err
	}

	err = db.ValidateAccountIds(ctx, account)
	if err != nil {
		db.Logger.Error(err, err.Error())
		return err
	}

	err = db.ValidateAccountParameters(ctx, account)
	if err != nil {
		db.Logger.Error(err, err.Error())
		return err
	}

	return nil
}

// ValidateAccountParameters validates account params.
func (db *Db) ValidateAccountParameters(ctx context.Context, account *merchant_service_proto_v1.MerchantAccount) error {
	err := db.ValidateAccountNotNil(ctx, account)
	if err != nil {
		db.Logger.Error(err, err.Error())
		return err
	}

	if account.BusinessEmail == constants.EMPTY || account.PhoneNumber == constants.EMPTY || account.BusinessName == constants.EMPTY {
		return errors.ErrMisconfiguredAccountParameters
	}

	return nil
}

// ValidateAccountNotNil ensures the account object is not nil
func (db *Db) ValidateAccountNotNil(ctx context.Context, account *merchant_service_proto_v1.MerchantAccount) error {
	if account != nil {
		return errors.ErrInvalidAccount
	}

	return nil
}

// ValidateAccountIds validates the existence of various ids associated with the account
func (db *Db) ValidateAccountIds(ctx context.Context, account *merchant_service_proto_v1.MerchantAccount) error {
	err := db.ValidateAccountNotNil(ctx, account)
	if err != nil {
		db.Logger.Error(err, err.Error())
		return err
	}

	if account.AuthnAccountId == 0 || account.PaymentsAccountId == 0 || account.StripeConnectedAccountId == 0 || account.EmployerId == 0 {
		return errors.ErrMisconfiguredIds
	}

	return nil
}

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

		if account.AccountState == merchant_service_proto_v1.MerchantAccountState_Inactive {
			return false, errors.ErrAccountDoesNotExist
		}

		return true, nil
	}
	return tx
}

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
	tx := func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
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

		if account.AccountState == merchant_service_proto_v1.MerchantAccountState_Inactive {
			return false, errors.ErrAccountDoesNotExist
		}

		return true, nil
	}
	return tx
}

// UpdateAccountOnboardStatus updates the onboarding status of a merchant account
func (db *Db) UpdateAccountOnboardStatus(ctx context.Context, account *merchant_service_proto_v1.MerchantAccount) error {
	err := db.ValidateAccountNotNil(ctx, account)
	if err != nil {
		db.Logger.Error(err, err.Error())
		return err
	}

	switch account.OnboardingState {
	// not started onboarding
	case merchant_service_proto_v1.OnboardingStatus_OnboardingNotStarted:
		account.OnboardingState = merchant_service_proto_v1.OnboardingStatus_FeelGuudOnboarding
		account.AccountState = merchant_service_proto_v1.MerchantAccountState_PendingOnboardingCompletion
		// completed onboarding with feelguud
	case merchant_service_proto_v1.OnboardingStatus_FeelGuudOnboarding:
		account.OnboardingState = merchant_service_proto_v1.OnboardingStatus_StripeOnboarding
		account.AccountState = merchant_service_proto_v1.MerchantAccountState_PendingOnboardingCompletion
		// completed onboarding with stripe
	case merchant_service_proto_v1.OnboardingStatus_StripeOnboarding:
		account.OnboardingState = merchant_service_proto_v1.OnboardingStatus_CatalogueOnboarding
		account.AccountState = merchant_service_proto_v1.MerchantAccountState_PendingOnboardingCompletion
		// completed onboarding catalogue
	case merchant_service_proto_v1.OnboardingStatus_CatalogueOnboarding:
		account.OnboardingState = merchant_service_proto_v1.OnboardingStatus_BCorpOnboarding
		account.AccountState = merchant_service_proto_v1.MerchantAccountState_ActiveAndOnboarded
	default:
		account.OnboardingState = merchant_service_proto_v1.OnboardingStatus_OnboardingNotStarted
		account.AccountState = merchant_service_proto_v1.MerchantAccountState_PendingOnboardingCompletion
	}

	return nil
}

// SaveAccountRecord saves a record in the database
func (db *Db) SaveAccountRecord(tx *gorm.DB, account *merchant_service_proto_v1.MerchantAccount) error {
	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(&account).Error; err != nil {
		return err
	}
	return nil
}
