package wrap

import "github.com/akm/sql-slog/public"

type (
	SqlLogger = public.SqlLogger
)

var (
	NewSqlLogger = public.NewSqlLogger

	IgnoreAttr  = public.IgnoreAttr
	WithNilAttr = public.WithNilAttr
)
