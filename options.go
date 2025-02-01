package sqlslog

import "github.com/akm/sql-slog/opts"

type (
	Option  = opts.Option
	Options = opts.Options
)

var (
	NewOptions = opts.NewOptions

	Logger                 = opts.Logger
	SetStepLogMsgFormatter = opts.SetStepLogMsgFormatter

	ConnBegin          = opts.ConnBegin
	ConnClose          = opts.ConnClose
	ConnPrepare        = opts.ConnPrepare
	ConnResetSession   = opts.ConnResetSession
	ConnPing           = opts.ConnPing
	ConnExecContext    = opts.ConnExecContext
	ConnQueryContext   = opts.ConnQueryContext
	ConnPrepareContext = opts.ConnPrepareContext
	ConnBeginTx        = opts.ConnBeginTx

	ConnectorConnect = opts.ConnectorConnect

	DriverOpen          = opts.DriverOpen
	DriverOpenConnector = opts.DriverOpenConnector

	SqlslogOpen = opts.SqlslogOpen

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
