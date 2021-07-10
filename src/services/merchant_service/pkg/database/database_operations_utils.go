package database

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/utils"
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

	if account.AuthnAccountId == 0 || account.PaymentsAccountId == 0 || account.StripeConnectedAccountId == constants.EMPTY || account.EmployerId == 0 {
		return errors.ErrMisconfiguredIds
	}

	return nil
}

func (db *Db) AccountActive(account *merchant_service_proto_v1.MerchantAccount) bool {
	if account == nil || !account.IsActive {
		return false
	}

	return true
}

// UpdateAccountOnboardStatus updates the onboarding status of a merchant account
func (db *Db) UpdateAccountOnboardStatus(ctx context.Context, account *merchant_service_proto_v1.MerchantAccount) error {
	err := db.ValidateAccountNotNil(ctx, account)
	if err != nil {
		db.Logger.Error(err, err.Error())
		return err
	}

	switch account.AccountOnboardingDetails {
	// not started onboarding
	case merchant_service_proto_v1.OnboardingStatus_OnboardingNotStarted:
		account.AccountOnboardingDetails = merchant_service_proto_v1.OnboardingStatus_FeelGuudOnboarding
		account.AccountOnboardingState = merchant_service_proto_v1.MerchantAccountState_PendingOnboardingCompletion
		// completed onboarding with feelguud
	case merchant_service_proto_v1.OnboardingStatus_FeelGuudOnboarding:
		account.AccountOnboardingDetails = merchant_service_proto_v1.OnboardingStatus_StripeOnboarding
		account.AccountOnboardingState = merchant_service_proto_v1.MerchantAccountState_PendingOnboardingCompletion
		// completed onboarding with stripe
	case merchant_service_proto_v1.OnboardingStatus_StripeOnboarding:
		account.AccountOnboardingDetails = merchant_service_proto_v1.OnboardingStatus_CatalogueOnboarding
		account.AccountOnboardingState = merchant_service_proto_v1.MerchantAccountState_PendingOnboardingCompletion
		// completed onboarding catalogue
	case merchant_service_proto_v1.OnboardingStatus_CatalogueOnboarding:
		account.AccountOnboardingDetails = merchant_service_proto_v1.OnboardingStatus_BCorpOnboarding
		account.AccountOnboardingState = merchant_service_proto_v1.MerchantAccountState_ActiveAndOnboarded
	default:
		account.AccountOnboardingDetails = merchant_service_proto_v1.OnboardingStatus_OnboardingNotStarted
		account.AccountOnboardingState = merchant_service_proto_v1.MerchantAccountState_PendingOnboardingCompletion
	}

	return nil
}

// startRootSpan starts a root span object
func (db *Db) startRootSpan(ctx context.Context, dbOpType OperationType) (context.Context, opentracing.Span) {
	return utils.StartRootOperationSpan(ctx, string(dbOpType), db.Logger)
}
