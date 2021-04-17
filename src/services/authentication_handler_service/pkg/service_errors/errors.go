package service_errors

import (
	"errors"
)

var (
	ErrInvalidRequest                           = errors.New("invalid grpc request")
	ErrFailedToConnectToDatabase                = errors.New("failed to connect to database")
	ErrFailedToPerformDatabaseMigrations        = errors.New("failed to perform database migrations")
	ErrInvalidInputArguments                    = errors.New("invalid input arguments")
	ErrInvalidEnvironmentVariableConfigurations = errors.New("invalid environment variable configurations")
	ErrFailedToStartGRPCServer                  = errors.New("failed to start grpc server")
	ErrHttpServerFailedGracefuleShutdown        = errors.New("http server failed to perform graceful shutdown")
	ErrHttpsServerFailedGracefuleShutdown       = errors.New("https server failed to perform graceful shutdown")
	ErrHttpServerCrashed                        = errors.New("http Server crashed")
	ErrHttpsServerCrashed                       = errors.New("https Server crashed")
	ErrSwaggerGenError                          = errors.New("swagger generation error")
	ErrFailedToWatchConfigDirectory             = errors.New("failed to watch config directory")
	ErrExceededMaxRetryAttempts                 = errors.New("exceeded max retry attemps")
	ErrInvalidAccount                           = errors.New("invalid account. account contains invalid fields")
	ErrFailedToReactivateAccount                = errors.New("failed to reactivate existing account")
	ErrDistributedTransactionError              = errors.New("distributed transaction error. failed to successfully perform a distributed operations")
	ErrFailedToUpdateAccountActiveStatus        = errors.New("failed to updated account active status")
	ErrAccountDoesNotExist                      = errors.New("account does not exist")
	ErrAccountAlreadyExist                      = errors.New("account already exists")
	ErrFailedToConvertFromOrmType               = errors.New("failed to perform conversion from Orm type")
	ErrFailedToConvertToOrmType                 = errors.New("failed to perform conversion to Orm type")
	ErrFailedToConfigureSaga                    = errors.New("failed to configure saga")
	ErrSagaFailedToExecuteSuccessfully          = errors.New("saga failed to execute successfully")
	ErrFailedToHashPassword                     = errors.New("failed to hash password")
	ErrFailedToCreateAccount                    = errors.New("failed to create account")
	ErrFailedToDeleteBusinessAccount            = errors.New("failed to delete business account")
	ErrFailedToUpdateAccountEmail               = errors.New("failed to updated account email through distributed transaction")
	ErrFailedToSaveUpdatedAccountRecord         = errors.New("failed to save updated account record")
	ErrCannotUpdatePassword                     = errors.New("cannot update password field")
	ErrCannotConfigureAccount                   = errors.New("cannot configure account")
	ErrUnableToObtainBusinessAccounts           = errors.New("unable to obtain account")
	ErrUnauthorizedRequest                      = errors.New("unauthorized request")
)

func NewError(msg string) error {
	return errors.New(msg)
}
