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
