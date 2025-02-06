package sqlslog

type Step string

func (s Step) String() string {
	return string(s)
}

const (
	StepConnBegin          Step = "Conn.Begin"
	StepConnBeginTx        Step = "Conn.BeginTx"
	StepConnClose          Step = "Conn.Close"
	StepConnPrepare        Step = "Conn.Prepare"
	StepConnPrepareContext Step = "Conn.PrepareContext"
	StepConnResetSession   Step = "Conn.ResetSession"
	StepConnPing           Step = "Conn.Ping"
	StepConnExecContext    Step = "Conn.ExecContext"
	StepConnQueryContext   Step = "Conn.QueryContext"

	StepConnectorConnect Step = "Connector.Connect"

	StepDriverOpen          Step = "Driver.Open"
	StepDriverOpenConnector Step = "Driver.OpenConnector"

	StepSqlslogOpen Step = "Open"

	StepRowsClose         Step = "Rows.Close"
	StepRowsNext          Step = "Rows.Next"
	StepRowsNextResultSet Step = "Rows.NextResultSet"

	StepStmtClose        Step = "Stmt.Close"
	StepStmtExec         Step = "Stmt.Exec"
	StepStmtQuery        Step = "Stmt.Query"
	StepStmtExecContext  Step = "Stmt.ExecContext"
	StepStmtQueryContext Step = "Stmt.QueryContext"

	StepTxCommit   Step = "Tx.Commit"
	StepTxRollback Step = "Tx.Rollback"
)
