package opts

type StmtOptions struct {
	Close        *StepOptions
	Exec         *StepOptions
	Query        *StepOptions
	ExecContext  *StepOptions
	QueryContext *StepOptions

	Rows *RowsOptions
}

func DefaultStmtOptions(formatter StepLogMsgFormatter) *StmtOptions {
	return &StmtOptions{
		Close:        DefaultStepOptions(formatter, "Stmt.Close", LevelInfo),
		Exec:         DefaultStepOptions(formatter, "Stmt.Exec", LevelInfo),
		Query:        DefaultStepOptions(formatter, "Stmt.Query", LevelInfo),
		ExecContext:  DefaultStepOptions(formatter, "Stmt.ExecContext", LevelInfo),
		QueryContext: DefaultStepOptions(formatter, "Stmt.QueryContext", LevelInfo),
		Rows:         DefaultRowsOptions(formatter),
	}
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
