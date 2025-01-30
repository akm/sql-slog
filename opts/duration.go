package opts

import (
	"log/slog"
	"time"
)

type DurationType int

const (
	DurationNanoSeconds  DurationType = iota // Duration in nanoseconds. Durations in log are expressed by slog.Int64
	DurationMicroSeconds                     // Duration in microseconds. Durations in log are expressed by slog.Int64
	DurationMilliSeconds                     // Duration in milliseconds. Durations in log are expressed by slog.Int64
	DurationGoDuration                       // Values in log are expressed with slog.Duration
	DurationString                           // Values in log are expressed with slog.String and time.Duration.String
)

func DurationAttrFunc(key string, t DurationType) func(d time.Duration) slog.Attr {
	switch t {
	case DurationNanoSeconds:
		return func(d time.Duration) slog.Attr { return slog.Int64(key, d.Nanoseconds()) }
	case DurationMicroSeconds:
		return func(d time.Duration) slog.Attr { return slog.Int64(key, d.Microseconds()) }
	case DurationMilliSeconds:
		return func(d time.Duration) slog.Attr { return slog.Int64(key, d.Milliseconds()) }
	case DurationGoDuration:
		return func(d time.Duration) slog.Attr { return slog.Duration(key, d) }
	case DurationString:
		return func(d time.Duration) slog.Attr { return slog.String(key, d.String()) }
	default:
		return func(d time.Duration) slog.Attr { return slog.Int64(key, d.Nanoseconds()) }
	}
}
