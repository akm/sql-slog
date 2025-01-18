package sqlslog

import "log/slog"

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

func newDefaultOptions(formatter StepLogMsgFormatter) *options {
	procOpts := func(name string, completeLevel Level) StepOptions {
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

	return &options{
		logger:              slog.Default(),
		connBegin:           procOpts("Conn.Begin", LevelInfo),
		connClose:           procOpts("Conn.Close", LevelInfo),
		connPrepare:         procOpts("Conn.Prepare", LevelInfo),
		connResetSession:    procOpts("Conn.ResetSession", LevelTrace),
		connPing:            procOpts("Conn.Ping", LevelTrace),
		connExecContext:     procOpts("Conn.ExecContext", LevelInfo),
		connQueryContext:    procOpts("Conn.QueryContext", LevelInfo),
		connPrepareContext:  procOpts("Conn.PrepareContext", LevelInfo),
		connBeginTx:         procOpts("Conn.BeginTx", LevelInfo),
		connectorConnect:    procOpts("Connector.Connect", LevelInfo),
		driverOpen:          procOpts("Driver.Open", LevelInfo),
		driverOpenConnector: procOpts("Driver.OpenConnector", LevelInfo),
		sqlslogOpen:         procOpts("sqlslog.Open", LevelInfo),
		rowsClose:           procOpts("Rows.Close", LevelDebug),
		rowsNext:            procOpts("Rows.Next", LevelDebug),
		rowsNextResultSet:   procOpts("Rows.NextResultSet", LevelDebug),
		stmtClose:           procOpts("Stmt.Close", LevelInfo),
		stmtExec:            procOpts("Stmt.Exec", LevelInfo),
		stmtQuery:           procOpts("Stmt.Query", LevelInfo),
		stmtExecContext:     procOpts("Stmt.ExecContext", LevelInfo),
		stmtQueryContext:    procOpts("Stmt.QueryContext", LevelInfo),
		txCommit:            procOpts("Tx.Commit", LevelInfo),
		txRollback:          procOpts("Tx.Rollback", LevelInfo),
	}
}

// Option is a function that sets some option on the options struct.
type Option func(*options)

var stepLogMsgFormatter = StepLogMsgWithEventName

// SetStepLogMsgFormatter sets the formatter for the step name used in logs.
// If not set, the default is StepLogMsgWithEventName.
func SetStepLogMsgFormatter(f StepLogMsgFormatter) { stepLogMsgFormatter = f }

func newOptions(opts ...Option) *options {
	o := newDefaultOptions(stepLogMsgFormatter)
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
