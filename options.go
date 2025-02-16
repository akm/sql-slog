package sqlslog

type options struct {
	stepLoggerOptions
	DriverOptions *driverOptions
	SlogOptions   *slogOptions
	Open          StepOptions
}

func newDefaultOptions(driverName string, msgb StepEventMsgBuilder) *options {
	return &options{
		stepLoggerOptions: stepLoggerOptions{
			durationKey:  DurationKeyDefault,
			durationType: DurationNanoSeconds,
		},
		DriverOptions: defaultDriverOptions(driverName, msgb),
		SlogOptions:   defaultSlogOptions(),
		Open:          *defaultStepOptions(msgb, StepSqlslogOpen, LevelInfo),
	}
}

// Option is a function that sets an option on the options struct.
type Option func(*options)

var stepEventMsgBuilder = StepEventMsgWithoutEventName

// SetStepEventMsgBuilder sets the builder for the step event message used in logs.
// If not set, the default is StepLogMsgWithEventName.
func SetStepEventMsgBuilder(f StepEventMsgBuilder) { stepEventMsgBuilder = f }

func newOptions(driverName string, opts ...Option) *options {
	o := newDefaultOptions(driverName, stepEventMsgBuilder)
	for _, opt := range opts {
		opt(o)
	}
	return o
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
