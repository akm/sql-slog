package testhelper

import (
	sqlslog "github.com/akm/sql-slog"
)

var StepEventMsgOptions = []sqlslog.Option{
	sqlslog.ConnBegin(func(o *sqlslog.StepOptions) { o.Start.Msg = "Conn.Begin Start" }),
	sqlslog.ConnBegin(func(o *sqlslog.StepOptions) { o.Error.Msg = "Conn.Begin Error" }),
	sqlslog.ConnBegin(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.Begin Complete" }),
	sqlslog.ConnClose(func(o *sqlslog.StepOptions) { o.Start.Msg = "Conn.Close Start" }),
	sqlslog.ConnClose(func(o *sqlslog.StepOptions) { o.Error.Msg = "Conn.Close Error" }),
	sqlslog.ConnClose(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.Close Complete" }),
	sqlslog.ConnPrepare(func(o *sqlslog.StepOptions) { o.Start.Msg = "Conn.Prepare Start" }),
	sqlslog.ConnPrepare(func(o *sqlslog.StepOptions) { o.Error.Msg = "Conn.Prepare Error" }),
	sqlslog.ConnPrepare(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.Prepare Complete" }),
	sqlslog.ConnResetSession(func(o *sqlslog.StepOptions) { o.Start.Msg = "Conn.ResetSession Start" }),
	sqlslog.ConnResetSession(func(o *sqlslog.StepOptions) { o.Error.Msg = "Conn.ResetSession Error" }),
	sqlslog.ConnResetSession(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.ResetSession Complete" }),
	sqlslog.ConnPing(func(o *sqlslog.StepOptions) { o.Start.Msg = "Conn.Ping Start" }),
	sqlslog.ConnPing(func(o *sqlslog.StepOptions) { o.Error.Msg = "Conn.Ping Error" }),
	sqlslog.ConnPing(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.Ping Complete" }),
	sqlslog.ConnExecContext(func(o *sqlslog.StepOptions) { o.Start.Msg = "Conn.ExecContext Start" }),
	sqlslog.ConnExecContext(func(o *sqlslog.StepOptions) { o.Error.Msg = "Conn.ExecContext Error" }),
	sqlslog.ConnExecContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.ExecContext Complete" }),
	sqlslog.ConnQueryContext(func(o *sqlslog.StepOptions) { o.Start.Msg = "Conn.QueryContext Start" }),
	sqlslog.ConnQueryContext(func(o *sqlslog.StepOptions) { o.Error.Msg = "Conn.QueryContext Error" }),
	sqlslog.ConnQueryContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.QueryContext Complete" }),
	sqlslog.ConnPrepareContext(func(o *sqlslog.StepOptions) { o.Start.Msg = "Conn.PrepareContext Start" }),
	sqlslog.ConnPrepareContext(func(o *sqlslog.StepOptions) { o.Error.Msg = "Conn.PrepareContext Error" }),
	sqlslog.ConnPrepareContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.PrepareContext Complete" }),
	sqlslog.ConnBeginTx(func(o *sqlslog.StepOptions) { o.Start.Msg = "Conn.BeginTx Start" }),
	sqlslog.ConnBeginTx(func(o *sqlslog.StepOptions) { o.Error.Msg = "Conn.BeginTx Error" }),
	sqlslog.ConnBeginTx(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Conn.BeginTx Complete" }),
	sqlslog.ConnectorConnect(func(o *sqlslog.StepOptions) { o.Start.Msg = "Connector.Connect Start" }),
	sqlslog.ConnectorConnect(func(o *sqlslog.StepOptions) { o.Error.Msg = "Connector.Connect Error" }),
	sqlslog.ConnectorConnect(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Connector.Connect Complete" }),
	sqlslog.DriverOpen(func(o *sqlslog.StepOptions) { o.Start.Msg = "Driver.Open Start" }),
	sqlslog.DriverOpen(func(o *sqlslog.StepOptions) { o.Error.Msg = "Driver.Open Error" }),
	sqlslog.DriverOpen(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Driver.Open Complete" }),
	sqlslog.DriverOpenConnector(func(o *sqlslog.StepOptions) { o.Start.Msg = "Driver.OpenConnector Start" }),
	sqlslog.DriverOpenConnector(func(o *sqlslog.StepOptions) { o.Error.Msg = "Driver.OpenConnector Error" }),
	sqlslog.DriverOpenConnector(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Driver.OpenConnector Complete" }),
	sqlslog.SqlslogOpen(func(o *sqlslog.StepOptions) { o.Start.Msg = "sqlslog.Open Start" }),
	sqlslog.SqlslogOpen(func(o *sqlslog.StepOptions) { o.Error.Msg = "sqlslog.Open Error" }),
	sqlslog.SqlslogOpen(func(o *sqlslog.StepOptions) { o.Complete.Msg = "sqlslog.Open Complete" }),
	sqlslog.RowsClose(func(o *sqlslog.StepOptions) { o.Start.Msg = "Rows.Close Start" }),
	sqlslog.RowsClose(func(o *sqlslog.StepOptions) { o.Error.Msg = "Rows.Close Error" }),
	sqlslog.RowsClose(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Rows.Close Complete" }),
	sqlslog.RowsNext(func(o *sqlslog.StepOptions) { o.Start.Msg = "Rows.Next Start" }),
	sqlslog.RowsNext(func(o *sqlslog.StepOptions) { o.Error.Msg = "Rows.Next Error" }),
	sqlslog.RowsNext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Rows.Next Complete" }),
	sqlslog.RowsNextResultSet(func(o *sqlslog.StepOptions) { o.Start.Msg = "Rows.NextResultSet Start" }),
	sqlslog.RowsNextResultSet(func(o *sqlslog.StepOptions) { o.Error.Msg = "Rows.NextResultSet Error" }),
	sqlslog.RowsNextResultSet(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Rows.NextResultSet Complete" }),
	sqlslog.StmtClose(func(o *sqlslog.StepOptions) { o.Start.Msg = "Stmt.Close Start" }),
	sqlslog.StmtClose(func(o *sqlslog.StepOptions) { o.Error.Msg = "Stmt.Close Error" }),
	sqlslog.StmtClose(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.Close Complete" }),
	sqlslog.StmtExec(func(o *sqlslog.StepOptions) { o.Start.Msg = "Stmt.Exec Start" }),
	sqlslog.StmtExec(func(o *sqlslog.StepOptions) { o.Error.Msg = "Stmt.Exec Error" }),
	sqlslog.StmtExec(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.Exec Complete" }),
	sqlslog.StmtQuery(func(o *sqlslog.StepOptions) { o.Start.Msg = "Stmt.Query Start" }),
	sqlslog.StmtQuery(func(o *sqlslog.StepOptions) { o.Error.Msg = "Stmt.Query Error" }),
	sqlslog.StmtQuery(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.Query Complete" }),
	sqlslog.StmtExecContext(func(o *sqlslog.StepOptions) { o.Start.Msg = "Stmt.ExecContext Start" }),
	sqlslog.StmtExecContext(func(o *sqlslog.StepOptions) { o.Error.Msg = "Stmt.ExecContext Error" }),
	sqlslog.StmtExecContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.ExecContext Complete" }),
	sqlslog.StmtQueryContext(func(o *sqlslog.StepOptions) { o.Start.Msg = "Stmt.QueryContext Start" }),
	sqlslog.StmtQueryContext(func(o *sqlslog.StepOptions) { o.Error.Msg = "Stmt.QueryContext Error" }),
	sqlslog.StmtQueryContext(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Stmt.QueryContext Complete" }),
	sqlslog.TxCommit(func(o *sqlslog.StepOptions) { o.Start.Msg = "Tx.Commit Start" }),
	sqlslog.TxCommit(func(o *sqlslog.StepOptions) { o.Error.Msg = "Tx.Commit Error" }),
	sqlslog.TxCommit(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Tx.Commit Complete" }),
	sqlslog.TxRollback(func(o *sqlslog.StepOptions) { o.Start.Msg = "Tx.Rollback Start" }),
	sqlslog.TxRollback(func(o *sqlslog.StepOptions) { o.Error.Msg = "Tx.Rollback Error" }),
	sqlslog.TxRollback(func(o *sqlslog.StepOptions) { o.Complete.Msg = "Tx.Rollback Complete" }),
}