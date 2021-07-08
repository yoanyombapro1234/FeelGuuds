package database

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	core_database "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-database"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-tracing"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/utils"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/saga"
)

type DatabaseOperationType string

// DbOperations provides an interface which any database tied to this service should implement
type DbOperations interface {
	CreateMerchantAccount(ctx context.Context, account *merchant_service_proto_v1.MerchantAccount) (*merchant_service_proto_v1.MerchantAccount, error)
	UpdateMerchantAccount(ctx context.Context, id uint64, account *merchant_service_proto_v1.MerchantAccount) (*merchant_service_proto_v1.MerchantAccount, error)
	DeleteMerchantAccount(ctx context.Context, id uint32) (bool, error)
	GetMerchantAccountById(ctx context.Context, id uint64) (*merchant_service_proto_v1.MerchantAccount, error)
	GetMerchantAccountsById(ctx context.Context, ids []uint32) ([]*merchant_service_proto_v1.MerchantAccount, error)
	DoesMerchantAccountExist(ctx context.Context, id uint64) (bool, error)
}

// Db witholds connection to a postgres database as well as a logging handler
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

// startRootSpan starts a root span object
func (db *Db) startRootSpan(ctx context.Context, dbOpType DatabaseOperationType) (context.Context, opentracing.Span) {
	return utils.StartRootOperationSpan(ctx, string(dbOpType), db.Logger)
}
