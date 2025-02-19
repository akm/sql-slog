package sqlslog

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
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

var stringToLevel = map[string]Level{
	"VERBOSE": LevelVerbose,
	"TRACE":   LevelTrace,
	"DEBUG":   LevelDebug,
	"INFO":    LevelInfo,
	"WARN":    LevelWarn,
	"ERROR":   LevelError,
}

var ErrUnknownLevel = errors.New("unknown level")

func ParseLevel(s string) (Level, error) {
	lv, ok := stringToLevel[strings.ToUpper(s)]
	if !ok {
		return 0, fmt.Errorf("%w: %q", ErrUnknownLevel, s)
	}
	return lv, nil
}

func ParseLevelWithDefault(s string, def Level) Level {
	lv, err := ParseLevel(s)
	if err != nil {
		return def
	}
	return lv
}

// ReplaceLevelAttr replaces the log level as sqlslog.Level with its string representation.
func ReplaceLevelAttr(_ []string, a slog.Attr) slog.Attr {
	// https://go.dev/src/log/slog/example_custom_levels_test.go
	if a.Key == slog.LevelKey {
		level := Level(a.Value.Any().(slog.Level)) //nolint:forcetypeassert
		a.Value = slog.StringValue(level.String())
	}
	return a
}
