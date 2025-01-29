package sqlslogopts

type DurationType int

const (
	DurationNanoSeconds  DurationType = iota // Duration in nanoseconds. Durations in log are expressed by slog.Int64
	DurationMicroSeconds                     // Duration in microseconds. Durations in log are expressed by slog.Int64
	DurationMilliSeconds                     // Duration in milliseconds. Durations in log are expressed by slog.Int64
	DurationGoDuration                       // Values in log are expressed with slog.Duration
	DurationString                           // Values in log are expressed with slog.String and time.Duration.String
)

// DurationKeyDefault is the default key for duration value in log.
const DurationKeyDefault = "duration"
