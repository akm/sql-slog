package sqlslog

import "github.com/akm/sql-slog/public"

type (
	Event               = public.Event
	StepOptions         = public.StepOptions
	StepLogMsgFormatter = public.StepLogMsgFormatter
)

var (
	NewStepOptions             = public.NewStepOptions
	StepLogMsgWithEventName    = public.StepLogMsgWithEventName
	StepLogMsgWithoutEventName = public.StepLogMsgWithoutEventName
)
