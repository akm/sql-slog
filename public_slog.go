package sqlslog

import (
	"github.com/akm/sql-slog/public"
)

type ReplaceAttrFunc = public.ReplaceAttrFunc

var (
	NewJSONHandler     = public.NewJSONHandler
	NewTextHandler     = public.NewTextHandler
	WrapHandlerOptions = public.WrapHandlerOptions
	MergeReplaceAttrs  = public.MergeReplaceAttrs
)
