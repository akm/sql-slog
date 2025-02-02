package sqlslog

import (
	"github.com/akm/sql-slog/internal/opts"
)

type (
	// Event is the event type of the step.
	Event = opts.Event

	// StepOptions is an struct that expresses the options for the step.
	StepOptions = opts.StepOptions
)

// StepLogMsgFormatter is the function type to format the step log message.
type StepLogMsgFormatter = func(name string, event Event) string

// SetStepLogMsgFormatter sets the formatter for the step name used in logs.
// If not set, the default is StepLogMsgWithEventName.
func SetStepLogMsgFormatter(f StepLogMsgFormatter) { opts.SetStepLogMsgFormatter(f) }
