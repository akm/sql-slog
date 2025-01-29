package wrap

import "github.com/akm/sql-slog/sqlslogopts"

type (
	SQLLogger = sqlslogopts.SQLLogger
)

var (
	NewSQLLogger = sqlslogopts.NewSQLLogger

	IgnoreAttr  = sqlslogopts.IgnoreAttr
	WithNilAttr = sqlslogopts.WithNilAttr
)
