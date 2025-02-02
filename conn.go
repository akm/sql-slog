package sqlslog

import "github.com/akm/sql-slog/internal/opts"

// Set the options for Conn.Begin.
func ConnBegin(f func(*StepOptions)) Option { return opts.ConnBegin(f) }

// Set the options for Conn.Close.
func ConnClose(f func(*StepOptions)) Option { return opts.ConnClose(f) }

// Set the options for Conn.Prepare.
func ConnPrepare(f func(*StepOptions)) Option { return opts.ConnPrepare(f) }

// Set the options for Conn.ResetSession.
func ConnResetSession(f func(*StepOptions)) Option { return opts.ConnResetSession(f) }

// Set the options for Conn.Ping.
func ConnPing(f func(*StepOptions)) Option { return opts.ConnPing(f) }

// Set the options for Conn.ExecContext.
func ConnExecContext(f func(*StepOptions)) Option { return opts.ConnExecContext(f) }

// Set the options for Conn.QueryContext.
func ConnQueryContext(f func(*StepOptions)) Option { return opts.ConnQueryContext(f) }

// Set the options for Conn.PrepareContext.
func ConnPrepareContext(f func(*StepOptions)) Option { return opts.ConnPrepareContext(f) }

// Set the options for Conn.BeginTx.
func ConnBeginTx(f func(*StepOptions)) Option { return opts.ConnBeginTx(f) }
