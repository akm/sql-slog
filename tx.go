package sqlslog

import "github.com/akm/sql-slog/internal/opts"

// Set the options for Tx.Commit.
func TxCommit(f func(*StepOptions)) Option { return opts.TxCommit(f) }

// Set the options for Tx.Rollback.
func TxRollback(f func(*StepOptions)) Option { return opts.TxRollback(f) }
