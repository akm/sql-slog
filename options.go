package sqlslog

import "github.com/akm/sql-slog/internal/opts"

type (
	// Option is a function that sets an option on the options struct.
	Option = opts.Option

	Options = opts.Options
)

var (
	RowsClose         = opts.RowsClose
	RowsNext          = opts.RowsNext
	RowsNextResultSet = opts.RowsNextResultSet

	StmtClose        = opts.StmtClose
	StmtExec         = opts.StmtExec
	StmtQuery        = opts.StmtQuery
	StmtExecContext  = opts.StmtExecContext
	StmtQueryContext = opts.StmtQueryContext

	TxCommit   = opts.TxCommit
	TxRollback = opts.TxRollback
)
