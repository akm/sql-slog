package sqlslog

import (
	"github.com/akm/sql-slog/opts"
)

type (
	Event       = opts.Event
	StepOptions = opts.StepOptions

	StepLogMsgFormatter = opts.StepLogMsgFormatter
)

var (
	NewStepOptions     = opts.NewStepOptions
	DefaultStepOptions = opts.DefaultStepOptions

	StepLogMsgWithEventName    = opts.StepLogMsgWithEventName
	StepLogMsgWithoutEventName = opts.StepLogMsgWithoutEventName
)
