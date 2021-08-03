package database

import (
	"context"
	"os"
	"time"

	core_database "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-database"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-tracing"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	svcErrors "github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/service_errors"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/utils"
	"gorm.io/gorm"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/saga"
)

type TxFunc func(ctx context.Context, tx *gorm.DB) (interface{}, error)

type OperationType string

// DbOperations provides an interface which any database tied to this service should implement
type DbOperations interface {
	CreateMerchantAccount(ctx context.Context, account *merchant_service_proto_v1.MerchantAccount) (*merchant_service_proto_v1.MerchantAccount, error)
	UpdateMerchantAccount(ctx context.Context, id uint64, account *merchant_service_proto_v1.MerchantAccount) (*merchant_service_proto_v1.MerchantAccount, error)
	DeactivateMerchantAccount(ctx context.Context, id uint64) (bool, error)
	GetMerchantAccountById(ctx context.Context, id uint64) (*merchant_service_proto_v1.MerchantAccount, error)
	GetMerchantAccountsById(ctx context.Context, ids []uint64) ([]*merchant_service_proto_v1.MerchantAccount, error)
	DoesMerchantAccountExist(ctx context.Context, id uint64) (bool, error)
	ActivateAccount(ctx context.Context, id uint64) (bool, error)
	UpdateAccountOnboardingStatus(ctx context.Context, id uint64, status merchant_service_proto_v1.MerchantAccountState) (*merchant_service_proto_v1.MerchantAccount, error)
}

// Db withholds connection to a postgres database as well as a logging handler
type Db struct {
	Conn                   *core_database.DatabaseConn
	Logger                 core_logging.ILog
	TracingEngine          *core_tracing.TracingEngine
	Saga                   *saga.SagaCoordinator
	MaxConnectionAttempts  int
	MaxRetriesPerOperation int
	RetryTimeOut           time.Duration
	OperationSleepInterval time.Duration
}

var _ DbOperations = (*Db)(nil)

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
