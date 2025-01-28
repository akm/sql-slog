package sqlslog

import "github.com/akm/sql-slog/public"

type (
	Options = public.Options
	Option  = public.Option
)

var (
	NewOptions             = public.NewOptions
	Logger                 = public.Logger
	SetStepLogMsgFormatter = public.SetStepLogMsgFormatter

	StmtQueryContext = public.StmtQueryContext
)
