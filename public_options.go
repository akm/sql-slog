package sqlslog

import "github.com/akm/sql-slog/public"

type (
	Options = public.Options
	Option  = public.Option
)

var (
	NewOptions             = public.NewOptions
	Logger                 = public.Logger
	SetStepLogMsgFormatter = public.SetStepLogMsgFormatter
)

var (
	ConnBegin           = public.ConnBegin
	ConnClose           = public.ConnClose
	ConnPrepare         = public.ConnPrepare
	ConnResetSession    = public.ConnResetSession
	ConnPing            = public.ConnPing
	ConnExecContext     = public.ConnExecContext
	ConnQueryContext    = public.ConnQueryContext
	ConnPrepareContext  = public.ConnPrepareContext
	ConnBeginTx         = public.ConnBeginTx
	ConnectorConnect    = public.ConnectorConnect
	DriverOpen          = public.DriverOpen
	DriverOpenConnector = public.DriverOpenConnector
	SqlslogOpen         = public.SqlslogOpen
	RowsClose           = public.RowsClose
	RowsNext            = public.RowsNext
	RowsNextResultSet   = public.RowsNextResultSet
	StmtClose           = public.StmtClose
	StmtExec            = public.StmtExec
	StmtQuery           = public.StmtQuery
	StmtExecContext     = public.StmtExecContext
	StmtQueryContext    = public.StmtQueryContext
	TxCommit            = public.TxCommit
	TxRollback          = public.TxRollback
)
