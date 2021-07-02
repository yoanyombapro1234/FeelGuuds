package service_errors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrRequestTimeout                           = status.Errorf(codes.DeadlineExceeded, "request timeout")
	ErrRetriesExceeded                          = status.Errorf(codes.DeadlineExceeded, "retries exceeded")
	ErrInvalidRequest                           = status.Errorf(codes.InvalidArgument, "invalid grpc request")
	ErrFailedToConnectToDatabase                = status.Errorf(codes.Internal, "failed to connect to database")
	ErrFailedToPerformDatabaseMigrations        = status.Errorf(codes.Internal, "failed to perform database migrations")
	ErrInvalidInputArguments                    = status.Errorf(codes.InvalidArgument, "invalid input arguments")
	ErrInvalidEnvironmentVariableConfigurations = status.Errorf(codes.Internal, "invalid environment variable configurations")
	ErrFailedToStartGRPCServer                  = status.Errorf(codes.Internal, "failed to start grpc server")
	ErrHttpServerFailedGracefuleShutdown        = status.Errorf(codes.Internal, "http server failed to perform graceful shutdown")
	ErrHttpsServerFailedGracefuleShutdown       = status.Errorf(codes.Internal, "https server failed to perform graceful shutdown")
	ErrHttpServerCrashed                        = status.Errorf(codes.Internal, "http Server crashed")
	ErrHttpsServerCrashed                       = status.Errorf(codes.Internal, "https Server crashed")
	ErrSwaggerGenError                          = status.Errorf(codes.Internal, "swagger generation error")
	ErrFailedToWatchConfigDirectory             = status.Errorf(codes.Internal, "failed to watch config directory")
	ErrExceededMaxRetryAttempts                 = status.Errorf(codes.DeadlineExceeded, "exceeded max retry attemps")
	ErrInvalidAccount                           = status.Errorf(codes.InvalidArgument, "invalid account. account contains invalid fields")
	ErrFailedToReactivateAccount                = status.Errorf(codes.FailedPrecondition, "failed to reactivate existing account")
	ErrDistributedTransactionError              = status.Errorf(codes.Internal,
		"distributed transaction error. failed to successfully perform a distributed operations")
	ErrFailedToUpdateAccountActiveStatus = status.Errorf(codes.Internal, "failed to updated account active status")
	ErrAccountDoesNotExist               = status.Errorf(codes.Internal, "account does not exist")
	ErrAccountAlreadyExist               = status.Errorf(codes.Internal, "account already exists")
	ErrFailedToConvertFromOrmType        = status.Errorf(codes.Internal, "failed to perform conversion from Orm type")
	ErrFailedToConvertToOrmType          = status.Errorf(codes.Internal, "failed to perform conversion to Orm type")
	ErrFailedToConfigureSaga             = status.Errorf(codes.Internal, "failed to configure saga")
	ErrSagaFailedToExecuteSuccessfully   = status.Errorf(codes.Internal, "saga failed to execute successfully")
	ErrFailedToHashPassword              = status.Errorf(codes.Internal, "failed to hash password")
	ErrFailedToCreateAccount             = status.Errorf(codes.Internal, "failed to create account")
	ErrFailedToUpdateAccountEmail        = status.Errorf(codes.Internal, "failed to updated account email through distributed transaction")
	ErrFailedToSaveUpdatedAccountRecord  = status.Errorf(codes.Internal, "failed to save updated account record")
	ErrCannotUpdatePassword              = status.Errorf(codes.Internal, "cannot update password field")
	ErrCannotConfigureAccount            = status.Errorf(codes.Internal, "cannot configure account")
	ErrUnableToObtainBusinessAccounts    = status.Errorf(codes.Internal, "unable to obtain account")
	ErrFailedToCastAccount               = status.Errorf(codes.Internal, "failed to cast account")
	ErrUnauthorizedRequest               = status.Errorf(codes.Unauthenticated, "unauthorized request")
)

// NewError returns a new error type based on some defined error message
func NewError(msg string) error {
	return errors.New(msg)
}
