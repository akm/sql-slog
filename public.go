package sqlslog

import (
	"log/slog"

	"github.com/akm/sql-slog/public"
)

type Level = public.Level

const (
	LevelVerbose Level = public.LevelVerbose
	LevelTrace   Level = public.LevelTrace
	LevelDebug   Level = public.LevelDebug
	LevelInfo    Level = public.LevelInfo
	LevelWarn    Level = public.LevelWarn
	LevelError   Level = public.LevelError
)

func ReplaceLevelAttr(_ []string, a slog.Attr) slog.Attr {
	return public.ReplaceLevelAttr(nil, a)
}
