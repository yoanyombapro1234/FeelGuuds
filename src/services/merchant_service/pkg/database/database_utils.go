package database

import (
	"context"
	"errors"
	"os"
	"time"

	core_database "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-database"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-tracing"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	svcErrors "github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/saga"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type ConnectionInitializationParams struct {
	ConnectionString       string
	TracingEngine          *core_tracing.TracingEngine
	Logger                 core_logging.ILog
	MaxConnectionAttempts  int
	MaxRetriesPerOperation int
	RetryTimeOut           time.Duration
	RetrySleepInterval     time.Duration
}

// New creates a database connection and returns the connection object
func New(ctx context.Context, params ConnectionInitializationParams) (*Db,
	error) {
	// generate a span for the database connection
	ctx, span := utils.StartRootOperationSpan(ctx, constants.DB_CONNECTION_ATTEMPT, params.Logger)
	defer span.Finish()

	if params.ConnectionString == constants.EMPTY || params.TracingEngine == nil || params.Logger == nil {
		// crash the process
		os.Exit(1)
	}

	logger := params.Logger

	dbConn := core_database.NewDatabaseConn(params.ConnectionString, "postgres")
	if dbConn == nil {
		logger.Fatal(svcErrors.ErrFailedToConnectToDatabase, svcErrors.ErrFailedToConnectToDatabase.Error())
	}
	logger.Info("Successfully connected to the database")

	ConfigureDatabaseConnection(dbConn)
	logger.Info("Successfully configured database connection object")

	err := MigrateSchemas(dbConn, logger, &merchant_service_proto_v1.MerchantAccount{})
	if err != nil {
		logger.Fatal(err, svcErrors.ErrFailedToPerformDatabaseMigrations.Error())
	}
	logger.Info("Successfully migrated database")

	return &Db{
		Conn:                   dbConn,
		Logger:                 logger,
		TracingEngine:          params.TracingEngine,
		Saga:                   saga.NewSagaCoordinator(logger),
		MaxConnectionAttempts:  params.MaxConnectionAttempts,
		MaxRetriesPerOperation: params.MaxRetriesPerOperation,
		RetryTimeOut:           params.RetryTimeOut,
		OperationSleepInterval: params.RetrySleepInterval,
	}, nil
}

// ConfigureDatabaseConnection configures a database connection
func ConfigureDatabaseConnection(dbConn *core_database.DatabaseConn) {
	dbConn.Engine.FullSaveAssociations = true
	dbConn.Engine.SkipDefaultTransaction = false
	dbConn.Engine.PrepareStmt = true
	dbConn.Engine.DisableAutomaticPing = false
	dbConn.Engine = dbConn.Engine.Set("gorm:auto_preload", true)
}

// MigrateSchemas creates or updates a given set of model based on a schema
// if it does not exist or migrates the model schemas to the latest version
func MigrateSchemas(db *core_database.DatabaseConn, logger core_logging.ILog, models ...interface{}) error {
	if err := db.Engine.AutoMigrate(models...); err != nil {
		// TODO: emit metric
		logger.Error(err, svcErrors.ErrFailedToPerformDatabaseMigrations.Error())
		return err
	}

	return nil
}

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
