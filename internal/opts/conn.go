package opts

import (
	"database/sql/driver"
	"errors"
	"log/slog"
)

type ConnOptions struct {
	IDGen IDGen

	Begin   *StepOptions
	BeginTx *StepOptions
	TxIDKey string
	Tx      *TxOptions

	Close *StepOptions

	Prepare        *StepOptions
	PrepareContext *StepOptions
	StmtIDKey      string
	Stmt           *StmtOptions

	ResetSession *StepOptions
	Ping         *StepOptions

	ExecContext  *StepOptions
	QueryContext *StepOptions
	Rows         *RowsOptions
}

const (
	TxIDKeyDefault   = "tx_id"
	StmtIDKeyDefault = "stmt_id"
)

func DefaultConnOptions(driverName string, formatter StepLogMsgFormatter) *ConnOptions {
	stmtOptions := DefaultStmtOptions(formatter)
	rowsOptions := stmtOptions.Rows

	return &ConnOptions{
		IDGen: IDGeneratorDefault,

		Begin:   DefaultStepOptions(formatter, "Conn.Begin", LevelInfo),
		BeginTx: DefaultStepOptions(formatter, "Conn.BeginTx", LevelInfo),
		TxIDKey: TxIDKeyDefault,
		Tx:      DefaultTxOptions(formatter),

		Close: DefaultStepOptions(formatter, "Conn.Close", LevelInfo),

		Prepare:        DefaultStepOptions(formatter, "Conn.Prepare", LevelInfo),
		PrepareContext: DefaultStepOptions(formatter, "Conn.PrepareContext", LevelInfo),
		StmtIDKey:      StmtIDKeyDefault,
		Stmt:           stmtOptions,

		ResetSession: DefaultStepOptions(formatter, "Conn.ResetSession", LevelTrace),
		Ping:         DefaultStepOptions(formatter, "Conn.Ping", LevelTrace),

		ExecContext:  DefaultStepOptions(formatter, "Conn.ExecContext", LevelInfo, ConnExecContextErrorHandler(driverName)),
		QueryContext: DefaultStepOptions(formatter, "Conn.QueryContext", LevelInfo, ConnQueryContextErrorHandler(driverName)),
		Rows:         rowsOptions,
	}
}

const (
	driverNameMysql = "mysql"
)

func ConnExecContextErrorHandler(driverName string) func(err error) (bool, []slog.Attr) {
	switch driverName {
	case driverNameMysql:
		return func(err error) (bool, []slog.Attr) {
			if err == nil {
				return true, nil
			}
			// https://pkg.go.dev/database/sql/driver#ErrSkip
			if errors.Is(err, driver.ErrSkip) {
				return true, []slog.Attr{slog.Bool("skip", true)}
			}
			return false, nil
		}
	default:
		return nil
	}
}

func ConnQueryContextErrorHandler(driverName string) func(err error) (bool, []slog.Attr) {
	switch driverName {
	case driverNameMysql:
		return func(err error) (bool, []slog.Attr) {
			if err == nil {
				return true, nil
			}
			// https://pkg.go.dev/database/sql/driver#ErrSkip
			if errors.Is(err, driver.ErrSkip) {
				return true, []slog.Attr{slog.Bool("skip", true)}
			}
			return false, nil
		}
	default:
		return nil
	}
}

// Set the options for Conn.Begin.
func ConnBegin(f func(*StepOptions)) Option { return func(o *Options) { f(o.Driver.Conn.Begin) } }

// Set the options for Conn.Close.
func ConnClose(f func(*StepOptions)) Option { return func(o *Options) { f(o.Driver.Conn.Close) } }

// Set the options for Conn.Prepare.
func ConnPrepare(f func(*StepOptions)) Option { return func(o *Options) { f(o.Driver.Conn.Prepare) } }

// Set the options for Conn.ResetSession.
func ConnResetSession(f func(*StepOptions)) Option {
	return func(o *Options) { f(o.Driver.Conn.ResetSession) }
}

// Set the options for Conn.Ping.
func ConnPing(f func(*StepOptions)) Option { return func(o *Options) { f(o.Driver.Conn.Ping) } }

// Set the options for Conn.ExecContext.
func ConnExecContext(f func(*StepOptions)) Option {
	return func(o *Options) { f(o.Driver.Conn.ExecContext) }
}

// Set the options for Conn.QueryContext.
func ConnQueryContext(f func(*StepOptions)) Option {
	return func(o *Options) { f(o.Driver.Conn.QueryContext) }
}

// Set the options for Conn.PrepareContext.
func ConnPrepareContext(f func(*StepOptions)) Option {
	return func(o *Options) { f(o.Driver.Conn.PrepareContext) }
}

// Set the options for Conn.BeginTx.
func ConnBeginTx(f func(*StepOptions)) Option { return func(o *Options) { f(o.Driver.Conn.BeginTx) } }
