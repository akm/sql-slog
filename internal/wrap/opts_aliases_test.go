package wrap

import "github.com/akm/sql-slog/internal/opts"

type DurationType = opts.DurationType

const (
	DurationNanoSeconds  = opts.DurationNanoSeconds
	DurationMicroSeconds = opts.DurationMicroSeconds
	DurationMilliSeconds = opts.DurationMilliSeconds
	DurationGoDuration   = opts.DurationGoDuration
	DurationString       = opts.DurationString
)

var NewJSONHandler = opts.NewJSONHandler
