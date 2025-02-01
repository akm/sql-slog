package sqlslog

import (
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
