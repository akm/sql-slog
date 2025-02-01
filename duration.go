package sqlslog

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

// Duration is an option to specify duration value in log.
// The default is DurationNanoSeconds.
func Duration(v DurationType) Option {
	return func(o *options) {
		o.durationType = v
	}
}

// DurationKey is an option to specify the key for duration value in log.
// The default is specified by DurationKeyDefault.
func DurationKey(key string) Option {
	return func(o *options) {
		o.durationKey = key
	}
}
