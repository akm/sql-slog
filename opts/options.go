package opts

import (
	"log/slog"
)

type Options struct {
	Logger *slog.Logger

	DurationKey  string
	DurationType DurationType

	IDGen     IDGen
	ConnIDKey string
	TxIDKey   string
	StmtIDKey string

	ConnBegin           StepOptions
	ConnClose           StepOptions
	ConnPrepare         StepOptions
	ConnResetSession    StepOptions
	ConnPing            StepOptions
	ConnExecContext     StepOptions
	ConnQueryContext    StepOptions
	ConnPrepareContext  StepOptions
	ConnBeginTx         StepOptions
	ConnectorConnect    StepOptions
	DriverOpen          StepOptions
	DriverOpenConnector StepOptions
	SqlslogOpen         StepOptions
	RowsClose           StepOptions
	RowsNext            StepOptions
	RowsNextResultSet   StepOptions
	StmtClose           StepOptions
	StmtExec            StepOptions
	StmtQuery           StepOptions
	StmtExecContext     StepOptions
	StmtQueryContext    StepOptions
	TxCommit            StepOptions
	TxRollback          StepOptions
}

func NewDefaultOptions(driverName string, formatter StepLogMsgFormatter) *Options {
	openOptions := DefaultOpenOptions(driverName, formatter)
	driverOptions := openOptions.Driver
	connectorOptions := driverOptions.Connector
	connOptions := connectorOptions.Conn
	stmtOptions := connOptions.Stmt
	rowsOptions := stmtOptions.Rows
	txOptions := connOptions.Tx

	return &Options{
		Logger:       slog.Default(),
		DurationKey:  DurationKeyDefault,
		DurationType: DurationNanoSeconds,

		IDGen:     connOptions.IDGen, // IDGeneratorDefault,
		ConnIDKey: driverOptions.ConnIDKey,
		TxIDKey:   connOptions.TxIDKey,
		StmtIDKey: connOptions.StmtIDKey,

		ConnBegin:           *connOptions.Begin,
		ConnClose:           *connOptions.Close,
		ConnPrepare:         *connOptions.Prepare,
		ConnResetSession:    *connOptions.ResetSession,
		ConnPing:            *connOptions.Ping,
		ConnExecContext:     *connOptions.ExecContext,
		ConnQueryContext:    *connOptions.QueryContext,
		ConnPrepareContext:  *connOptions.PrepareContext,
		ConnBeginTx:         *connOptions.BeginTx,
		ConnectorConnect:    *connectorOptions.Connect,
		DriverOpen:          *driverOptions.Open,
		DriverOpenConnector: *driverOptions.OpenConnector,
		SqlslogOpen:         *openOptions.Open,
		RowsClose:           *rowsOptions.Close,
		RowsNext:            *rowsOptions.Next,
		RowsNextResultSet:   *rowsOptions.NextResultSet,
		StmtClose:           *stmtOptions.Close,
		StmtExec:            *stmtOptions.Exec,
		StmtQuery:           *stmtOptions.Query,
		StmtExecContext:     *stmtOptions.ExecContext,
		StmtQueryContext:    *stmtOptions.QueryContext,
		TxCommit:            *txOptions.Commit,
		TxRollback:          *txOptions.Rollback,
	}
}

// DurationKeyDefault is the default key for duration value in log.
const DurationKeyDefault = "duration"

// Option is a function that sets an option on the options struct.
type Option func(*Options)

var stepLogMsgFormatter = StepLogMsgWithoutEventName

// SetStepLogMsgFormatter sets the formatter for the step name used in logs.
// If not set, the default is StepLogMsgWithEventName.
func SetStepLogMsgFormatter(f StepLogMsgFormatter) { stepLogMsgFormatter = f }

func NewOptions(driverName string, opts ...Option) *Options {
	o := NewDefaultOptions(driverName, stepLogMsgFormatter)
	for _, opt := range opts {
		opt(o)
	}
	return o
}
