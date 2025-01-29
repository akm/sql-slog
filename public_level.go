package sqlslog

import (
	"log/slog"

	sqlslogopts "github.com/akm/sql-slog/opts"
)

type Level = sqlslogopts.Level

const (
	LevelVerbose Level = sqlslogopts.LevelVerbose
	LevelTrace   Level = sqlslogopts.LevelTrace
	LevelDebug   Level = sqlslogopts.LevelDebug
	LevelInfo    Level = sqlslogopts.LevelInfo
	LevelWarn    Level = sqlslogopts.LevelWarn
	LevelError   Level = sqlslogopts.LevelError
)

func ReplaceLevelAttr(_ []string, a slog.Attr) slog.Attr {
	return sqlslogopts.ReplaceLevelAttr(nil, a)
}
