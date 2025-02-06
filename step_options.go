package sqlslog

import "log/slog"

type EventOptions struct {
	Msg   string
	Level Level
}

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

// StepLogMsgFormatter is the function type to format the step log message.
type StepLogMsgFormatter func(step Step, event Event) string

// StepLogMsgWithEventName returns the formatted step log message with the event name.
func StepLogMsgWithEventName(step Step, event Event) string {
	return step.String() + " " + event.String()
}

// StepLogMsgWithoutEventName returns the formatted step log message without the event name.
func StepLogMsgWithoutEventName(step Step, _ Event) string {
	return step.String()
}

// StepOptions is an struct that expresses the options for the step.
type StepOptions struct {
	Start    EventOptions
	Error    EventOptions
	Complete EventOptions

	// ErrorHandler is the function to handle the error.
	// When the error should not be logged as an error but as complete, it should return true.
	// It can also add attributes to the log.
	ErrorHandler func(error) (bool, []slog.Attr)
}

const defaultSlogLevelDiff = 4

func (o *StepOptions) SetLevel(lv Level) {
	o.Start.Level = lv - defaultSlogLevelDiff
	o.Complete.Level = lv
}

func (o *StepOptions) compare(other *StepOptions) bool {
	return o.Start.Level == other.Start.Level &&
		o.Error.Level == other.Error.Level &&
		o.Complete.Level == other.Complete.Level
}

func newStepOptions(f StepLogMsgFormatter, step Step, startLevel, errorLevel, completeLevel Level) *StepOptions {
	return &StepOptions{
		Start:    EventOptions{Msg: f(step, EventStart), Level: startLevel},
		Error:    EventOptions{Msg: f(step, EventError), Level: errorLevel},
		Complete: EventOptions{Msg: f(step, EventComplete), Level: completeLevel},
	}
}

func defaultStepOptions(formatter StepLogMsgFormatter, step Step, completeLevel Level, errHandlers ...func(error) (bool, []slog.Attr)) *StepOptions { // nolint:unparam
	var startLevel Level
	switch completeLevel { // nolint:exhaustive
	case LevelError:
		startLevel = LevelInfo
	case LevelInfo:
		startLevel = LevelDebug
	case LevelDebug:
		startLevel = LevelTrace
	default:
		startLevel = LevelVerbose
	}
	r := newStepOptions(formatter, step, startLevel, LevelError, completeLevel)
	if len(errHandlers) > 0 {
		r.ErrorHandler = errHandlers[0]
	}
	return r
}
