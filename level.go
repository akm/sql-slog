package sqlslog

import (
	"github.com/akm/sql-slog/internal/opts"
)

// Level is the log level for sqlslog.
type Level = opts.Level

const (
	LevelVerbose Level = opts.LevelVerbose // Lower than slog.LevelTrace.
	LevelTrace   Level = opts.LevelTrace   // Lower than slog.LevelDebug.
	LevelDebug   Level = opts.LevelDebug   // Same as slog.LevelDebug.
	LevelInfo    Level = opts.LevelInfo    // Same as slog.LevelInfo.
	LevelWarn    Level = opts.LevelWarn    // Same as slog.LevelWarn.
	LevelError   Level = opts.LevelError   // Same as slog.LevelError.
)
