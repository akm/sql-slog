package testhelper

import (
	sqlslog "github.com/akm/sql-slog"
)

var StepEventMsgOptions = []sqlslog.Option{
	sqlslog.ConnBegin(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.Begin Complete" }),
	sqlslog.ConnClose(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.Close Complete" }),
	sqlslog.ConnPrepare(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.Prepare Complete" }),
	sqlslog.ConnResetSession(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.ResetSession Complete" }),
	sqlslog.ConnPing(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.Ping Complete" }),
	sqlslog.ConnExecContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.ExecContext Complete" }),
	sqlslog.ConnQueryContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.QueryContext Complete" }),
	sqlslog.ConnPrepareContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.PrepareContext Complete" }),
	sqlslog.ConnBeginTx(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.BeginTx Complete" }),
	sqlslog.ConnectorConnect(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Connector.Connect Complete" }),
	sqlslog.DriverOpen(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Driver.Open Complete" }),
	sqlslog.DriverOpenConnector(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Driver.OpenConnector Complete" }),
	sqlslog.SqlslogOpen(func(o *sqlslog.StepOptions) { o.Complete.Msg = "sqlslog.Open Complete" }),
	sqlslog.RowsClose(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Rows.Close Complete" }),
	sqlslog.RowsNext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Rows.Next Complete" }),
	sqlslog.RowsNextResultSet(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Rows.NextResultSet Complete" }),
	sqlslog.StmtClose(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.Close Complete" }),
	sqlslog.StmtExec(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.Exec Complete" }),
	sqlslog.StmtQuery(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.Query Complete" }),
	sqlslog.StmtExecContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.ExecContext Complete" }),
	sqlslog.StmtQueryContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.QueryContext Complete" }),
	sqlslog.TxCommit(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Tx.Commit Complete" }),
	sqlslog.TxRollback(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Tx.Rollback Complete" }),
}
