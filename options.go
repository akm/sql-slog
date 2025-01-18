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

type Option func(*options)

var procNameFormatter = ProcNameWithEventName

func SetProcNameFormatter(f ProcNameFormatter) { procNameFormatter = f }

func newOptions(opts ...Option) *options {
	o := newDefaultOptions(procNameFormatter)
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func Logger(logger *slog.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// options for Conn
func ConnBegin(opts *ProcOptions) Option   { return func(o *options) { o.connBegin = *opts } }
func ConnPrepare(opts *ProcOptions) Option { return func(o *options) { o.connPrepare = *opts } }
func ConnResetSession(opts *ProcOptions) Option {
	return func(o *options) { o.connResetSession = *opts }
}
func ConnPing(opts *ProcOptions) Option { return func(o *options) { o.connPing = *opts } }
func ConnExecContext(opts *ProcOptions) Option {
	return func(o *options) { o.connExecContext = *opts }
}
func ConnQueryContext(opts *ProcOptions) Option {
	return func(o *options) { o.connQueryContext = *opts }
}
func ConnPrepareContext(opts *ProcOptions) Option {
	return func(o *options) { o.connPrepareContext = *opts }
}
func ConnBeginTx(opts *ProcOptions) Option { return func(o *options) { o.connBeginTx = *opts } }

// options for Connector
func ConnectorConnect(opts *ProcOptions) Option {
	return func(o *options) { o.connectorConnect = *opts }
}

// options for Driver
func DriverOpen(opts *ProcOptions) Option { return func(o *options) { o.driverOpen = *opts } }
func DriverOpenConnector(opts *ProcOptions) Option {
	return func(o *options) { o.driverOpenConnector = *opts }
}

// options for sqlslog
func SqlslogOpen(opts *ProcOptions) Option { return func(o *options) { o.sqlslogOpen = *opts } }

// options for Rows
func RowsClose(opts *ProcOptions) Option { return func(o *options) { o.rowsClose = *opts } }
func RowsNext(opts *ProcOptions) Option  { return func(o *options) { o.rowsNext = *opts } }
func RowsNextResultSet(opts *ProcOptions) Option {
	return func(o *options) { o.rowsNextResultSet = *opts }
}

// options for Stmt
func StmtClose(opts *ProcOptions) Option { return func(o *options) { o.stmtClose = *opts } }
func StmtExec(opts *ProcOptions) Option  { return func(o *options) { o.stmtExec = *opts } }
func StmtQuery(opts *ProcOptions) Option { return func(o *options) { o.stmtQuery = *opts } }
func StmtExecContext(opts *ProcOptions) Option {
	return func(o *options) { o.stmtExecContext = *opts }
}
func StmtQueryContext(opts *ProcOptions) Option {
	return func(o *options) { o.stmtQueryContext = *opts }
}

// options for Tx
func TxCommit(opts *ProcOptions) Option   { return func(o *options) { o.txCommit = *opts } }
func TxRollback(opts *ProcOptions) Option { return func(o *options) { o.txRollback = *opts } }
