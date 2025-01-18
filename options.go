package sqlslog

import "log/slog"

type options struct {
	logger *slog.Logger

	connBegin           ProcOptions
	connClose           ProcOptions
	connPrepare         ProcOptions
	connResetSession    ProcOptions
	connPing            ProcOptions
	connExecContext     ProcOptions
	connQueryContext    ProcOptions
	connPrepareContext  ProcOptions
	connBeginTx         ProcOptions
	connectorConnect    ProcOptions
	driverOpen          ProcOptions
	driverOpenConnector ProcOptions
	sqlslogOpen         ProcOptions
	rowsClose           ProcOptions
	rowsNext            ProcOptions
	rowsNextResultSet   ProcOptions
	stmtClose           ProcOptions
	stmtExec            ProcOptions
	stmtQuery           ProcOptions
	stmtExecContext     ProcOptions
	stmtQueryContext    ProcOptions
	txCommit            ProcOptions
	txRollback          ProcOptions
}

func newDefaultOptions(formatter ProcNameFormatter) *options {
	procOpts := func(name string, completeLevel Level) ProcOptions {
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
		return *newProcOptions(formatter, name, startLevel, LevelError, completeLevel)
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

var procNameFormatter = ProcNameWithEventName

// SetProcNameFormatter sets the formatter for the process name used in logs.
// If not set, the default is ProcNameWithEventName.
func SetProcNameFormatter(f ProcNameFormatter) { procNameFormatter = f }

func newOptions(opts ...Option) *options {
	o := newDefaultOptions(procNameFormatter)
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
func ConnBegin(opts *ProcOptions) Option { return func(o *options) { o.connBegin = *opts } }

// Set the options for Conn.Prepare
func ConnPrepare(opts *ProcOptions) Option { return func(o *options) { o.connPrepare = *opts } }

// Set the options for Conn.ResetSession
func ConnResetSession(opts *ProcOptions) Option {
	return func(o *options) { o.connResetSession = *opts }
}

// Set the options for Conn.Ping
func ConnPing(opts *ProcOptions) Option { return func(o *options) { o.connPing = *opts } }

// Set the options for Conn.ExecContext
func ConnExecContext(opts *ProcOptions) Option {
	return func(o *options) { o.connExecContext = *opts }
}

// Set the options for Conn.QueryContext
func ConnQueryContext(opts *ProcOptions) Option {
	return func(o *options) { o.connQueryContext = *opts }
}

// Set the options for Conn.PrepareContext
func ConnPrepareContext(opts *ProcOptions) Option {
	return func(o *options) { o.connPrepareContext = *opts }
}

// Set the options for Conn.BeginTx
func ConnBeginTx(opts *ProcOptions) Option { return func(o *options) { o.connBeginTx = *opts } }

// Set the options for Connector.Connect
func ConnectorConnect(opts *ProcOptions) Option {
	return func(o *options) { o.connectorConnect = *opts }
}

// Set the options for Driver.Open
func DriverOpen(opts *ProcOptions) Option { return func(o *options) { o.driverOpen = *opts } }

// Set the options for Driver.OpenConnector
func DriverOpenConnector(opts *ProcOptions) Option {
	return func(o *options) { o.driverOpenConnector = *opts }
}

// Set the options for sqlslog.Open
func SqlslogOpen(opts *ProcOptions) Option { return func(o *options) { o.sqlslogOpen = *opts } }

// Set the options for Rows.Close
func RowsClose(opts *ProcOptions) Option { return func(o *options) { o.rowsClose = *opts } }

// Set the options for Rows.Next
func RowsNext(opts *ProcOptions) Option { return func(o *options) { o.rowsNext = *opts } }

// Set the options for Rows.NextResultSet
func RowsNextResultSet(opts *ProcOptions) Option {
	return func(o *options) { o.rowsNextResultSet = *opts }
}

// Set the options for Stmt.Close
func StmtClose(opts *ProcOptions) Option { return func(o *options) { o.stmtClose = *opts } }

// Set the options for Stmt.Exec
func StmtExec(opts *ProcOptions) Option { return func(o *options) { o.stmtExec = *opts } }

// Set the options for Stmt.Query
func StmtQuery(opts *ProcOptions) Option { return func(o *options) { o.stmtQuery = *opts } }

// Set the options for Stmt.ExecContext
func StmtExecContext(opts *ProcOptions) Option {
	return func(o *options) { o.stmtExecContext = *opts }
}

// Set the options for Stmt.QueryContext
func StmtQueryContext(opts *ProcOptions) Option {
	return func(o *options) { o.stmtQueryContext = *opts }
}

// Set the options for Tx.Commit
func TxCommit(opts *ProcOptions) Option { return func(o *options) { o.txCommit = *opts } }

// Set the options for Tx.Rollback
func TxRollback(opts *ProcOptions) Option { return func(o *options) { o.txRollback = *opts } }
