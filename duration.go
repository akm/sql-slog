package sqlslog

import "github.com/akm/sql-slog/internal/opts"

type DurationType = opts.DurationType

const (
	DurationNanoSeconds  = opts.DurationNanoSeconds
	DurationMicroSeconds = opts.DurationMicroSeconds
	DurationMilliSeconds = opts.DurationMilliSeconds
	DurationGoDuration   = opts.DurationGoDuration
	DurationString       = opts.DurationString
)

// Duration is an option to specify duration value in log.
// The default is DurationNanoSeconds.
func Duration(v DurationType) Option { return opts.Duration(v) }

// DurationKey is an option to specify the key for duration value in log.
// The default is specified by DurationKeyDefault.
func DurationKey(key string) Option { return opts.DurationKey(key) }
