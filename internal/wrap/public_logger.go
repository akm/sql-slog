package wrap

import "github.com/akm/sql-slog/public"

type (
	SQLLogger = public.SQLLogger
)

var (
	NewSqlLogger = public.NewSqlLogger

	IgnoreAttr  = public.IgnoreAttr
	WithNilAttr = public.WithNilAttr
)
