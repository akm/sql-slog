package wrap

import sqlslogopts "github.com/akm/sql-slog/opts"

type (
	SQLLogger = sqlslogopts.SQLLogger
)

var (
	NewSQLLogger = sqlslogopts.NewSQLLogger

	IgnoreAttr  = sqlslogopts.IgnoreAttr
	WithNilAttr = sqlslogopts.WithNilAttr
)
