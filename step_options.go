package sqlslog

import (
	"log/slog"

	"github.com/akm/sql-slog/opts"
)

type (
	Event       = opts.Event
	StepOptions = opts.StepOptions

	StepLogMsgFormatter = opts.StepLogMsgFormatter
)

var (
	NewStepOptions = opts.NewStepOptions

	StepLogMsgWithEventName    = opts.StepLogMsgWithEventName
	StepLogMsgWithoutEventName = opts.StepLogMsgWithoutEventName
)

func DefaultStepOptions(formatter StepLogMsgFormatter, name string, completeLevel Level, errHandlers ...func(error) (bool, []slog.Attr)) *StepOptions {
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
	r := NewStepOptions(formatter, name, startLevel, LevelError, completeLevel)
	if len(errHandlers) > 0 {
		r.ErrorHandler = errHandlers[0]
	}
	return r
}
