package sqlslog

import "github.com/akm/sql-slog/sqlslogopts"

type (
	Event               = sqlslogopts.Event
	StepOptions         = sqlslogopts.StepOptions
	StepLogMsgFormatter = sqlslogopts.StepLogMsgFormatter
)

var (
	NewStepOptions             = sqlslogopts.NewStepOptions
	StepLogMsgWithEventName    = sqlslogopts.StepLogMsgWithEventName
	StepLogMsgWithoutEventName = sqlslogopts.StepLogMsgWithoutEventName
)
