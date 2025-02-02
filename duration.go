package sqlslog

import "github.com/akm/sql-slog/internal/opts"

type DurationType = opts.DurationType

const (
	DurationNanoSeconds  = opts.DurationNanoSeconds  // Duration in nanoseconds. Durations in log are expressed by slog.Int64
	DurationMicroSeconds = opts.DurationMicroSeconds // Duration in microseconds. Durations in log are expressed by slog.Int64
	DurationMilliSeconds = opts.DurationMilliSeconds // Duration in milliseconds. Durations in log are expressed by slog.Int64
	DurationGoDuration   = opts.DurationGoDuration   // Values in log are expressed with slog.Duration
	DurationString       = opts.DurationString       // Values in log are expressed with slog.String and time.Duration.String
)

// Duration is an option to specify duration value in log.
// The default is DurationNanoSeconds.
func Duration(v DurationType) Option { return opts.Duration(v) }

// DurationKey is an option to specify the key for duration value in log.
// The default is specified by DurationKeyDefault.
func DurationKey(key string) Option { return opts.DurationKey(key) }
