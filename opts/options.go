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

// Set the options for Conn.Begin.
func ConnBegin(f func(*StepOptions)) Option { return func(o *Options) { f(&o.ConnBegin) } }

// Set the options for Conn.Close.
func ConnClose(f func(*StepOptions)) Option { return func(o *Options) { f(&o.ConnClose) } }

// Set the options for Conn.Prepare.
func ConnPrepare(f func(*StepOptions)) Option { return func(o *Options) { f(&o.ConnPrepare) } }

// Set the options for Conn.ResetSession.
func ConnResetSession(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.ConnResetSession) }
}

// Set the options for Conn.Ping.
func ConnPing(f func(*StepOptions)) Option { return func(o *Options) { f(&o.ConnPing) } }

// Set the options for Conn.ExecContext.
func ConnExecContext(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.ConnExecContext) }
}

// Set the options for Conn.QueryContext.
func ConnQueryContext(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.ConnQueryContext) }
}

// Set the options for Conn.PrepareContext.
func ConnPrepareContext(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.ConnPrepareContext) }
}

// Set the options for Conn.BeginTx.
func ConnBeginTx(f func(*StepOptions)) Option { return func(o *Options) { f(&o.ConnBeginTx) } }

// Set the options for Connector.Connect.
func ConnectorConnect(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.ConnectorConnect) }
}

// Set the options for Driver.Open.
func DriverOpen(f func(*StepOptions)) Option { return func(o *Options) { f(&o.DriverOpen) } }

// Set the options for Driver.OpenConnector.
func DriverOpenConnector(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.DriverOpenConnector) }
}

// Set the options for sqlslog.Open.
func SqlslogOpen(f func(*StepOptions)) Option { return func(o *Options) { f(&o.SqlslogOpen) } } // nolint:revive

// Set the options for Rows.Close.
func RowsClose(f func(*StepOptions)) Option { return func(o *Options) { f(&o.RowsClose) } }

// Set the options for Rows.Next.
func RowsNext(f func(*StepOptions)) Option { return func(o *Options) { f(&o.RowsNext) } }

// Set the options for Rows.NextResultSet.
func RowsNextResultSet(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.RowsNextResultSet) }
}

// Set the options for Stmt.Close.
func StmtClose(f func(*StepOptions)) Option { return func(o *Options) { f(&o.StmtClose) } }

// Set the options for Stmt.Exec.
func StmtExec(f func(*StepOptions)) Option { return func(o *Options) { f(&o.StmtExec) } }

// Set the options for Stmt.Query.
func StmtQuery(f func(*StepOptions)) Option { return func(o *Options) { f(&o.StmtQuery) } }

// Set the options for Stmt.ExecContext.
func StmtExecContext(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.StmtExecContext) }
}

// Set the options for Stmt.QueryContext.
func StmtQueryContext(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.StmtQueryContext) }
}

// Set the options for Tx.Commit.
func TxCommit(f func(*StepOptions)) Option { return func(o *Options) { f(&o.TxCommit) } }

// Set the options for Tx.Rollback.
func TxRollback(f func(*StepOptions)) Option { return func(o *Options) { f(&o.TxRollback) } }
