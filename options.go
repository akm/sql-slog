package sqlslog

import (
	"log/slog"
)

type options struct {
	logger *slog.Logger

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

func newDefaultOptions(driverName string, formatter StepLogMsgFormatter) *options {
	stepOpts := func(name string, completeLevel Level) StepOptions {
		var startLevel Level
		switch completeLevel {
		case LevelError:
			startLevel = LevelInfo
		case LevelInfo:
			startLevel = LevelDebug
		case LevelDebug:
			startLevel = LevelTrace
		default:
			startLevel = LevelVerbose
		}
		return *newStepOptions(formatter, name, startLevel, LevelError, completeLevel)
	}

	withErrorHandler := func(opts StepOptions, eh func(error) (bool, []slog.Attr)) StepOptions {
		r := opts
		r.ErrorHandler = eh
		return r
	}

	return &options{
		logger:              slog.Default(),
		connBegin:           stepOpts("Conn.Begin", LevelInfo),
		connClose:           stepOpts("Conn.Close", LevelInfo),
		connPrepare:         stepOpts("Conn.Prepare", LevelInfo),
		connResetSession:    stepOpts("Conn.ResetSession", LevelTrace),
		connPing:            stepOpts("Conn.Ping", LevelTrace),
		connExecContext:     stepOpts("Conn.ExecContext", LevelInfo),
		connQueryContext:    stepOpts("Conn.QueryContext", LevelInfo),
		connPrepareContext:  stepOpts("Conn.PrepareContext", LevelInfo),
		connBeginTx:         stepOpts("Conn.BeginTx", LevelInfo),
		connectorConnect:    withErrorHandler(stepOpts("Connector.Connect", LevelInfo), ConnectorConnectErrorHandler(driverName)),
		driverOpen:          withErrorHandler(stepOpts("Driver.Open", LevelInfo), HandleDriverOpenError),
		driverOpenConnector: stepOpts("Driver.OpenConnector", LevelInfo),
		sqlslogOpen:         stepOpts("sqlslog.Open", LevelInfo),
		rowsClose:           stepOpts("Rows.Close", LevelDebug),
		rowsNext:            withErrorHandler(stepOpts("Rows.Next", LevelDebug), HandleRowsNextError),
		rowsNextResultSet:   stepOpts("Rows.NextResultSet", LevelDebug),
		stmtClose:           stepOpts("Stmt.Close", LevelInfo),
		stmtExec:            stepOpts("Stmt.Exec", LevelInfo),
		stmtQuery:           stepOpts("Stmt.Query", LevelInfo),
		stmtExecContext:     stepOpts("Stmt.ExecContext", LevelInfo),
		stmtQueryContext:    stepOpts("Stmt.QueryContext", LevelInfo),
		txCommit:            stepOpts("Tx.Commit", LevelInfo),
		txRollback:          stepOpts("Tx.Rollback", LevelInfo),
	}
}

// Option is a function that sets some option on the options struct.
type Option func(*options)

var stepLogMsgFormatter = StepLogMsgWithEventName

// SetStepLogMsgFormatter sets the formatter for the step name used in logs.
// If not set, the default is StepLogMsgWithEventName.
func SetStepLogMsgFormatter(f StepLogMsgFormatter) { stepLogMsgFormatter = f }

func newOptions(driverName string, opts ...Option) *options {
	o := newDefaultOptions(driverName, stepLogMsgFormatter)
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Logger sets the slog.Logger to be used.
// If not set, the default is slog.Default().
func Logger(logger *slog.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// Set the options for Conn.Begin
func ConnBegin(f func(*StepOptions)) Option { return func(o *options) { f(&o.connBegin) } }

// Set the options for Conn.Close
func ConnClose(f func(*StepOptions)) Option { return func(o *options) { f(&o.connClose) } }

// Set the options for Conn.Prepare
func ConnPrepare(f func(*StepOptions)) Option { return func(o *options) { f(&o.connPrepare) } }

// Set the options for Conn.ResetSession
func ConnResetSession(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.connResetSession) }
}

// Set the options for Conn.Ping
func ConnPing(f func(*StepOptions)) Option { return func(o *options) { f(&o.connPing) } }

// Set the options for Conn.ExecContext
func ConnExecContext(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.connExecContext) }
}

// Set the options for Conn.QueryContext
func ConnQueryContext(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.connQueryContext) }
}

// Set the options for Conn.PrepareContext
func ConnPrepareContext(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.connPrepareContext) }
}

// Set the options for Conn.BeginTx
func ConnBeginTx(f func(*StepOptions)) Option { return func(o *options) { f(&o.connBeginTx) } }

// Set the options for Connector.Connect
func ConnectorConnect(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.connectorConnect) }
}

// Set the options for Driver.Open
func DriverOpen(f func(*StepOptions)) Option { return func(o *options) { f(&o.driverOpen) } }

// Set the options for Driver.OpenConnector
func DriverOpenConnector(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.driverOpenConnector) }
}

// Set the options for sqlslog.Open
func SqlslogOpen(f func(*StepOptions)) Option { return func(o *options) { f(&o.sqlslogOpen) } }

// Set the options for Rows.Close
func RowsClose(f func(*StepOptions)) Option { return func(o *options) { f(&o.rowsClose) } }

// Set the options for Rows.Next
func RowsNext(f func(*StepOptions)) Option { return func(o *options) { f(&o.rowsNext) } }

// Set the options for Rows.NextResultSet
func RowsNextResultSet(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.rowsNextResultSet) }
}

// Set the options for Stmt.Close
func StmtClose(f func(*StepOptions)) Option { return func(o *options) { f(&o.stmtClose) } }

// Set the options for Stmt.Exec
func StmtExec(f func(*StepOptions)) Option { return func(o *options) { f(&o.stmtExec) } }

// Set the options for Stmt.Query
func StmtQuery(f func(*StepOptions)) Option { return func(o *options) { f(&o.stmtQuery) } }

// Set the options for Stmt.ExecContext
func StmtExecContext(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.stmtExecContext) }
}

// Set the options for Stmt.QueryContext
func StmtQueryContext(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.stmtQueryContext) }
}

// Set the options for Tx.Commit
func TxCommit(f func(*StepOptions)) Option { return func(o *options) { f(&o.txCommit) } }

// Set the options for Tx.Rollback
func TxRollback(f func(*StepOptions)) Option { return func(o *options) { f(&o.txRollback) } }
