package sqlslog

import (
	sqlslogopts "github.com/akm/sql-slog/opts"
)

type ReplaceAttrFunc = sqlslogopts.ReplaceAttrFunc

var (
	NewJSONHandler     = sqlslogopts.NewJSONHandler
	NewTextHandler     = sqlslogopts.NewTextHandler
	WrapHandlerOptions = sqlslogopts.WrapHandlerOptions
	MergeReplaceAttrs  = sqlslogopts.MergeReplaceAttrs
)
