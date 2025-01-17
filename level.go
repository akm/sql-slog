package sqlslog

import (
	"fmt"
	"log/slog"
)

type Level slog.Level

const (
	LevelVerbose Level = Level(-12)
	LevelTrace   Level = Level(-8)
	LevelDebug   Level = Level(slog.LevelDebug)
	LevelInfo    Level = Level(slog.LevelInfo)
	LevelWarn    Level = Level(slog.LevelWarn)
	LevelError   Level = Level(slog.LevelError)
)

// var _ slog.Level = slog.Level(LevelVerbose)

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

func ReplaceLevelAttr(_ []string, a slog.Attr) slog.Attr {
	// https://go.dev/src/log/slog/example_custom_levels_test.go
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(Level)
		a.Value = slog.StringValue(level.String())
	}
	return a
}
