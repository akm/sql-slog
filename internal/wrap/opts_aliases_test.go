package wrap

import "github.com/akm/sql-slog/opts"

type DurationType = opts.DurationType

const (
	DurationNanoSeconds  = opts.DurationNanoSeconds
	DurationMicroSeconds = opts.DurationMicroSeconds
	DurationMilliSeconds = opts.DurationMilliSeconds
	DurationGoDuration   = opts.DurationGoDuration
	DurationString       = opts.DurationString
)

var DurationAttrFunc = opts.DurationAttrFunc
