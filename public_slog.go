package sqlslog

import (
	"github.com/akm/sql-slog/sqlslogopts"
)

type ReplaceAttrFunc = sqlslogopts.ReplaceAttrFunc

var (
	NewJSONHandler     = sqlslogopts.NewJSONHandler
	NewTextHandler     = sqlslogopts.NewTextHandler
	WrapHandlerOptions = sqlslogopts.WrapHandlerOptions
	MergeReplaceAttrs  = sqlslogopts.MergeReplaceAttrs
)
