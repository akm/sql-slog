package sqlslog

import "github.com/akm/sql-slog/internal/opts"

// Set the options for Stmt.Close.
func StmtClose(f func(*StepOptions)) Option { return opts.StmtClose(f) }

// Set the options for Stmt.Exec.
func StmtExec(f func(*StepOptions)) Option { return opts.StmtExec(f) }

// Set the options for Stmt.Query.
func StmtQuery(f func(*StepOptions)) Option { return opts.StmtQuery(f) }

// Set the options for Stmt.ExecContext.
func StmtExecContext(f func(*StepOptions)) Option { return opts.StmtExecContext(f) }

// Set the options for Stmt.QueryContext.
func StmtQueryContext(f func(*StepOptions)) Option { return opts.StmtQueryContext(f) }
