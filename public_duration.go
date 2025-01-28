package sqlslog

import (
	"github.com/akm/sql-slog/public"
)

type DurationType = public.DurationType

const (
	DurationNanoSeconds  DurationType = public.DurationNanoSeconds
	DurationMicroSeconds DurationType = public.DurationMicroSeconds
	DurationMilliSeconds DurationType = public.DurationMilliSeconds
	DurationGoDuration   DurationType = public.DurationGoDuration
	DurationString       DurationType = public.DurationString
)

const DurationKeyDefault = public.DurationKeyDefault
