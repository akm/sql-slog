package wrap

import "github.com/akm/sql-slog/public"

type (
	SQLLogger = public.SQLLogger
)

var (
	NewSQLLogger = public.NewSQLLogger

	IgnoreAttr  = public.IgnoreAttr
	WithNilAttr = public.WithNilAttr
)
