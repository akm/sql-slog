package sqlslog

import (
	"log/slog"

	"github.com/akm/sql-slog/public"
)

type Options struct {
	Logger *slog.Logger

	DurationKey  string
	DurationType DurationType

	IdGen     IDGen
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
	stepOpts := func(name string, completeLevel Level) StepOptions {
		var startLevel Level
		switch completeLevel { // nolint:exhaustive
		case LevelError:
			startLevel = LevelInfo
		case LevelInfo:
			startLevel = LevelDebug
		case LevelDebug:
			startLevel = LevelTrace
		default:
			startLevel = LevelVerbose
		}
		return *NewStepOptions(formatter, name, startLevel, LevelError, completeLevel)
	}

	withErrorHandler := func(opts StepOptions, eh func(error) (bool, []slog.Attr)) StepOptions {
		r := opts
		r.ErrorHandler = eh
		return r
	}

	return &Options{
		Logger:       slog.Default(),
		DurationKey:  DurationKeyDefault,
		DurationType: DurationNanoSeconds,

		IdGen:     IDGeneratorDefault,
		ConnIDKey: ConnIDKeyDefault,
		TxIDKey:   TxIDKeyDefault,
		StmtIDKey: StmtIDKeyDefault,

		ConnBegin:           stepOpts("Conn.Begin", LevelInfo),
		ConnClose:           stepOpts("Conn.Close", LevelInfo),
		ConnPrepare:         stepOpts("Conn.Prepare", LevelInfo),
		ConnResetSession:    stepOpts("Conn.ResetSession", LevelTrace),
		ConnPing:            stepOpts("Conn.Ping", LevelTrace),
		ConnExecContext:     withErrorHandler(stepOpts("Conn.ExecContext", LevelInfo), public.ConnExecContextErrorHandler(driverName)),
		ConnQueryContext:    withErrorHandler(stepOpts("Conn.QueryContext", LevelInfo), public.ConnQueryContextErrorHandler(driverName)),
		ConnPrepareContext:  stepOpts("Conn.PrepareContext", LevelInfo),
		ConnBeginTx:         stepOpts("Conn.BeginTx", LevelInfo),
		ConnectorConnect:    withErrorHandler(stepOpts("Connector.Connect", LevelInfo), public.ConnectorConnectErrorHandler(driverName)),
		DriverOpen:          withErrorHandler(stepOpts("Driver.Open", LevelInfo), public.DriverOpenErrorHandler(driverName)),
		DriverOpenConnector: stepOpts("Driver.OpenConnector", LevelInfo),
		SqlslogOpen:         stepOpts("sqlslog.Open", LevelInfo),
		RowsClose:           stepOpts("Rows.Close", LevelDebug),
		RowsNext:            withErrorHandler(stepOpts("Rows.Next", LevelDebug), public.HandleRowsNextError),
		RowsNextResultSet:   stepOpts("Rows.NextResultSet", LevelDebug),
		StmtClose:           stepOpts("Stmt.Close", LevelInfo),
		StmtExec:            stepOpts("Stmt.Exec", LevelInfo),
		StmtQuery:           stepOpts("Stmt.Query", LevelInfo),
		StmtExecContext:     stepOpts("Stmt.ExecContext", LevelInfo),
		StmtQueryContext:    stepOpts("Stmt.QueryContext", LevelInfo),
		TxCommit:            stepOpts("Tx.Commit", LevelInfo),
		TxRollback:          stepOpts("Tx.Rollback", LevelInfo),
	}
}

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
		o.Logger = logger
	}
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
