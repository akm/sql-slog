package sqlslog

import (
	"log/slog"

	"github.com/akm/sql-slog/opts"
)

type Level = opts.Level

const (
	LevelVerbose Level = opts.LevelVerbose
	LevelTrace   Level = opts.LevelTrace
	LevelDebug   Level = opts.LevelDebug
	LevelInfo    Level = opts.LevelInfo
	LevelWarn    Level = opts.LevelWarn
	LevelError   Level = opts.LevelError
)

// ReplaceLevelAttr replaces the log level as sqlslog.Level with its string representation.
func ReplaceLevelAttr(_ []string, a slog.Attr) slog.Attr {
	// https://go.dev/src/log/slog/example_custom_levels_test.go
	if a.Key == slog.LevelKey {
		level := opts.Level(a.Value.Any().(slog.Level)) //nolint:forcetypeassert
		a.Value = slog.StringValue(level.String())
	}
	return a
}
