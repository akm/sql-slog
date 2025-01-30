package opts

import (
	"fmt"
	"log/slog"
)

// Level is the log level for sqlslog.
type Level slog.Level

const (
	LevelVerbose Level = Level(-12)             // Lower than slog.LevelTrace.
	LevelTrace   Level = Level(-8)              // Lower than slog.LevelDebug.
	LevelDebug   Level = Level(slog.LevelDebug) // Same as slog.LevelDebug.
	LevelInfo    Level = Level(slog.LevelInfo)  // Same as slog.LevelInfo.
	LevelWarn    Level = Level(slog.LevelWarn)  // Same as slog.LevelWarn.
	LevelError   Level = Level(slog.LevelError) // Same as slog.LevelError.
)

var _ slog.Leveler = LevelVerbose

// String returns the string representation of the log level.
func (l Level) String() string {
	str := func(base string, val Level) string {
		if val == 0 {
			return base
		}
		return fmt.Sprintf("%s%+d", base, val)
	}

	switch {
	case l < LevelTrace:
		return str("VERBOSE", l-LevelVerbose)
	case l < LevelDebug:
		return str("TRACE", l-LevelTrace)
	case l < LevelInfo:
		return str("DEBUG", l-LevelDebug)
	case l < LevelWarn:
		return str("INFO", l-LevelInfo)
	case l < LevelError:
		return str("WARN", l-LevelWarn)
	default:
		return str("ERROR", l-LevelError)
	}
}

// Level returns the slog.Level.
func (l Level) Level() slog.Level {
	return slog.Level(l)
}
