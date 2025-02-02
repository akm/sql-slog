package sqlslog

import "github.com/akm/sql-slog/internal/opts"

// Set the options for Rows.Close.
func RowsClose(f func(*StepOptions)) Option { return opts.RowsClose(f) }

// Set the options for Rows.Next.
func RowsNext(f func(*StepOptions)) Option { return opts.RowsNext(f) }

// Set the options for Rows.NextResultSet.
func RowsNextResultSet(f func(*StepOptions)) Option { return opts.RowsNextResultSet(f) }
