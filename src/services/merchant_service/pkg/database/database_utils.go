package database

import (
	"context"
	"errors"

	svcErr "github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/service_errors"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/utils"

	"github.com/opentracing/opentracing-go"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	"golang.org/x/crypto/bcrypt"
)

// ValidateAndHashPassword validates, hashes and salts a password
func (db *Db) ValidateAndHashPassword(password string) (string, error) {
	// check if confirmed password is not empty
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	//  hash and salt password
	hashedPassword, err := db.hashAndSalt([]byte(password))
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

// hashAndSalt hashes and salts a password
func (db *Db) hashAndSalt(pwd []byte) (string, error) {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

// ComparePasswords compares a hashed password and a plaintext password and returns
// a boolean stating wether they are equal or not
func (db *Db) ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}
	return true
}

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
		return svcErr.ErrMisconfiguredAccountParameters
	}

	return nil
}

// ValidateAccountNotNil ensures the account object is not nil
func (db *Db) ValidateAccountNotNil(ctx context.Context, account *merchant_service_proto_v1.MerchantAccount) error {
	if account != nil {
		return svcErr.ErrInvalidAccount
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
		return svcErr.ErrMisconfiguredIds
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
