package sqlslog

import (
	"log/slog"
)

type options struct {
	stepLoggerOptions
	sqlslogOptions
}

func newDefaultOptions(driverName string, formatter StepLogMsgFormatter) *options {
	return &options{
		stepLoggerOptions: stepLoggerOptions{
			logger:       slog.Default(),
			durationKey:  DurationKeyDefault,
			durationType: DurationNanoSeconds,
		},
		sqlslogOptions: *defaultSqlslogOptions(driverName, formatter),
	}
}

// Option is a function that sets an option on the options struct.
type Option func(*options)

var stepLogMsgFormatter = StepLogMsgWithoutEventName

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
		o.stepLoggerOptions.logger = logger
	}
}

// Set the options for Conn.Begin.
func ConnBegin(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.Begin) }
}

// Set the options for Conn.Close.
func ConnClose(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.Close) }
}

// Set the options for Conn.Prepare.
func ConnPrepare(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.Prepare) }
}

// Set the options for Conn.ResetSession.
func ConnResetSession(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.ResetSession) }
}

// Set the options for Conn.Ping.
func ConnPing(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.Ping) }
}

// Set the options for Conn.ExecContext.
func ConnExecContext(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.ExecContext) }
}

// Set the options for Conn.QueryContext.
func ConnQueryContext(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.QueryContext) }
}

// Set the options for Conn.PrepareContext.
func ConnPrepareContext(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.PrepareContext) }
}

// Set the options for Conn.BeginTx.
func ConnBeginTx(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.BeginTx) }
}

// Set the options for Connector.Connect.
func ConnectorConnect(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnectorOptions.Connect) }
}

// Set the options for Driver.Open.
func DriverOpen(f func(*StepOptions)) Option { return func(o *options) { f(&o.DriverOptions.Open) } }

// Set the options for Driver.OpenConnector.
func DriverOpenConnector(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.OpenConnector) }
}

// Set the options for sqlslog.Open.
func SqlslogOpen(f func(*StepOptions)) Option { return func(o *options) { f(&o.Open) } } // nolint:revive

// Set the options for Rows.Close.
func RowsClose(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.RowsOptions.Close) }
}

// Set the options for Rows.Next.
func RowsNext(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.RowsOptions.Next) }
}

// Set the options for Rows.NextResultSet.
func RowsNextResultSet(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.RowsOptions.NextResultSet) }
}

// Set the options for Stmt.Close.
func StmtClose(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.StmtOptions.Close) }
}

// Set the options for Stmt.Exec.
func StmtExec(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.StmtOptions.Exec) }
}

// Set the options for Stmt.Query.
func StmtQuery(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.StmtOptions.Query) }
}

// Set the options for Stmt.ExecContext.
func StmtExecContext(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.StmtOptions.ExecContext) }
}

// Set the options for Stmt.QueryContext.
func StmtQueryContext(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.StmtOptions.QueryContext) }
}

// Set the options for Tx.Commit.
func TxCommit(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.TxOptions.Commit) }
}

// Set the options for Tx.Rollback.
func TxRollback(f func(*StepOptions)) Option {
	return func(o *options) { f(&o.DriverOptions.ConnOptions.TxOptions.Rollback) }
}
