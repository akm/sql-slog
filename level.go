package sqlslog

import (
	"fmt"
	"log/slog"
)

// Level is log level for sqlslog.
type Level slog.Level

const (
	LevelVerbose Level = Level(-12)             // lower than slog.LevelTrace.
	LevelTrace   Level = Level(-8)              // lower than slog.LevelDebug.
	LevelDebug   Level = Level(slog.LevelDebug) // the same as slog.LevelDebug.
	LevelInfo    Level = Level(slog.LevelInfo)  // the same as slog.LevelInfo.
	LevelWarn    Level = Level(slog.LevelWarn)  // the same as slog.LevelWarn.
	LevelError   Level = Level(slog.LevelError) // the same as slog.LevelError.
)

// var _ slog.Level = slog.Level(LevelVerbose)
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

// ReplaceLevelAttr replaces the log level as sqlslog.Level with the string representation.
func ReplaceLevelAttr(_ []string, a slog.Attr) slog.Attr {
	// https://go.dev/src/log/slog/example_custom_levels_test.go
	if a.Key == slog.LevelKey {
		level := Level(a.Value.Any().(slog.Level))
		a.Value = slog.StringValue(level.String())
	}
	return a
}
