package sqlslog

import (
	"log/slog"
)

type Options struct {
	logger *slog.Logger

	durationKey  string
	durationType DurationType

	idGen     IDGen
	connIDKey string
	txIDKey   string
	stmtIDKey string

	connBegin           StepOptions
	connClose           StepOptions
	connPrepare         StepOptions
	connResetSession    StepOptions
	connPing            StepOptions
	connExecContext     StepOptions
	connQueryContext    StepOptions
	connPrepareContext  StepOptions
	connBeginTx         StepOptions
	connectorConnect    StepOptions
	driverOpen          StepOptions
	driverOpenConnector StepOptions
	sqlslogOpen         StepOptions
	rowsClose           StepOptions
	rowsNext            StepOptions
	rowsNextResultSet   StepOptions
	stmtClose           StepOptions
	stmtExec            StepOptions
	stmtQuery           StepOptions
	stmtExecContext     StepOptions
	stmtQueryContext    StepOptions
	txCommit            StepOptions
	txRollback          StepOptions
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
		logger:       slog.Default(),
		durationKey:  DurationKeyDefault,
		durationType: DurationNanoSeconds,

		idGen:     connOptions.idGen, // IDGeneratorDefault,
		connIDKey: driverOptions.connIDKey,
		txIDKey:   connOptions.txIDKey,
		stmtIDKey: connOptions.stmtIDKey,

		connBegin:           *connOptions.Begin,
		connClose:           *connOptions.Close,
		connPrepare:         *connOptions.Prepare,
		connResetSession:    *connOptions.ResetSession,
		connPing:            *connOptions.Ping,
		connExecContext:     *connOptions.ExecContext,
		connQueryContext:    *connOptions.QueryContext,
		connPrepareContext:  *connOptions.PrepareContext,
		connBeginTx:         *connOptions.BeginTx,
		connectorConnect:    *connectorOptions.Connect,
		driverOpen:          *driverOptions.Open,
		driverOpenConnector: *driverOptions.OpenConnector,
		sqlslogOpen:         *openOptions.Open,
		rowsClose:           *rowsOptions.Close,
		rowsNext:            *rowsOptions.Next,
		rowsNextResultSet:   *rowsOptions.NextResultSet,
		stmtClose:           *stmtOptions.Close,
		stmtExec:            *stmtOptions.Exec,
		stmtQuery:           *stmtOptions.Query,
		stmtExecContext:     *stmtOptions.ExecContext,
		stmtQueryContext:    *stmtOptions.QueryContext,
		txCommit:            *txOptions.Commit,
		txRollback:          *txOptions.Rollback,
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

// Logger sets the slog.Logger to be used.
// If not set, the default is slog.Default().
func Logger(logger *slog.Logger) Option {
	return func(o *Options) {
		o.logger = logger
	}
}

// Set the options for Conn.Begin.
func ConnBegin(f func(*StepOptions)) Option { return func(o *Options) { f(&o.connBegin) } }

// Set the options for Conn.Close.
func ConnClose(f func(*StepOptions)) Option { return func(o *Options) { f(&o.connClose) } }

// Set the options for Conn.Prepare.
func ConnPrepare(f func(*StepOptions)) Option { return func(o *Options) { f(&o.connPrepare) } }

// Set the options for Conn.ResetSession.
func ConnResetSession(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.connResetSession) }
}

// Set the options for Conn.Ping.
func ConnPing(f func(*StepOptions)) Option { return func(o *Options) { f(&o.connPing) } }

// Set the options for Conn.ExecContext.
func ConnExecContext(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.connExecContext) }
}

// Set the options for Conn.QueryContext.
func ConnQueryContext(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.connQueryContext) }
}

// Set the options for Conn.PrepareContext.
func ConnPrepareContext(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.connPrepareContext) }
}

// Set the options for Conn.BeginTx.
func ConnBeginTx(f func(*StepOptions)) Option { return func(o *Options) { f(&o.connBeginTx) } }

// Set the options for Connector.Connect.
func ConnectorConnect(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.connectorConnect) }
}

// Set the options for Driver.Open.
func DriverOpen(f func(*StepOptions)) Option { return func(o *Options) { f(&o.driverOpen) } }

// Set the options for Driver.OpenConnector.
func DriverOpenConnector(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.driverOpenConnector) }
}

// Set the options for sqlslog.Open.
func SqlslogOpen(f func(*StepOptions)) Option { return func(o *Options) { f(&o.sqlslogOpen) } } // nolint:revive

// Set the options for Rows.Close.
func RowsClose(f func(*StepOptions)) Option { return func(o *Options) { f(&o.rowsClose) } }

// Set the options for Rows.Next.
func RowsNext(f func(*StepOptions)) Option { return func(o *Options) { f(&o.rowsNext) } }

// Set the options for Rows.NextResultSet.
func RowsNextResultSet(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.rowsNextResultSet) }
}

// Set the options for Stmt.Close.
func StmtClose(f func(*StepOptions)) Option { return func(o *Options) { f(&o.stmtClose) } }

// Set the options for Stmt.Exec.
func StmtExec(f func(*StepOptions)) Option { return func(o *Options) { f(&o.stmtExec) } }

// Set the options for Stmt.Query.
func StmtQuery(f func(*StepOptions)) Option { return func(o *Options) { f(&o.stmtQuery) } }

// Set the options for Stmt.ExecContext.
func StmtExecContext(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.stmtExecContext) }
}

// Set the options for Stmt.QueryContext.
func StmtQueryContext(f func(*StepOptions)) Option {
	return func(o *Options) { f(&o.stmtQueryContext) }
}

// Set the options for Tx.Commit.
func TxCommit(f func(*StepOptions)) Option { return func(o *Options) { f(&o.txCommit) } }

// Set the options for Tx.Rollback.
func TxRollback(f func(*StepOptions)) Option { return func(o *Options) { f(&o.txRollback) } }
