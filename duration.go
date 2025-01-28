package sqlslog

type DurationType int

const (
	DurationNanoSeconds  DurationType = iota // Duration in nanoseconds. Durations in log are expressed by slog.Int64
	DurationMicroSeconds                     // Duration in microseconds. Durations in log are expressed by slog.Int64
	DurationMilliSeconds                     // Duration in milliseconds. Durations in log are expressed by slog.Int64
	DurationGoDuration                       // Values in log are expressed with slog.Duration
	DurationString                           // Values in log are expressed with slog.String and time.Duration.String
)

// Duration is an option to specify duration value in log.
// The default is DurationNanoSeconds.
func Duration(v DurationType) Option {
	return func(o *Options) {
		o.durationType = v
	}
}

// DurationKey is an option to specify the key for duration value in log.
// The default is specified by DurationKeyDefault.
func DurationKey(key string) Option {
	return func(o *Options) {
		o.durationKey = key
	}
}

// DurationKeyDefault is the default key for duration value in log.
const DurationKeyDefault = "duration"
