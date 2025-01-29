package sqlslog

import sqlslogopts "github.com/akm/sql-slog/opts"

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
