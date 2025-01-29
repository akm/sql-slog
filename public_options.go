package sqlslog

import "github.com/akm/sql-slog/sqlslogopts"

type (
	Options = sqlslogopts.Options
	Option  = sqlslogopts.Option
)

var (
	NewOptions             = sqlslogopts.NewOptions
	Logger                 = sqlslogopts.Logger
	SetStepLogMsgFormatter = sqlslogopts.SetStepLogMsgFormatter
)

var (
	ConnBegin           = sqlslogopts.ConnBegin
	ConnClose           = sqlslogopts.ConnClose
	ConnPrepare         = sqlslogopts.ConnPrepare
	ConnResetSession    = sqlslogopts.ConnResetSession
	ConnPing            = sqlslogopts.ConnPing
	ConnExecContext     = sqlslogopts.ConnExecContext
	ConnQueryContext    = sqlslogopts.ConnQueryContext
	ConnPrepareContext  = sqlslogopts.ConnPrepareContext
	ConnBeginTx         = sqlslogopts.ConnBeginTx
	ConnectorConnect    = sqlslogopts.ConnectorConnect
	DriverOpen          = sqlslogopts.DriverOpen
	DriverOpenConnector = sqlslogopts.DriverOpenConnector
	SqlslogOpen         = sqlslogopts.SqlslogOpen
	RowsClose           = sqlslogopts.RowsClose
	RowsNext            = sqlslogopts.RowsNext
	RowsNextResultSet   = sqlslogopts.RowsNextResultSet
	StmtClose           = sqlslogopts.StmtClose
	StmtExec            = sqlslogopts.StmtExec
	StmtQuery           = sqlslogopts.StmtQuery
	StmtExecContext     = sqlslogopts.StmtExecContext
	StmtQueryContext    = sqlslogopts.StmtQueryContext
	TxCommit            = sqlslogopts.TxCommit
	TxRollback          = sqlslogopts.TxRollback
)
