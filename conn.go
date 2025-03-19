package sqlslog

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"log/slog"
)

type connOptions struct {
	IDGen IDGen

	Begin     StepOptions
	BeginTx   StepOptions
	TxIDKey   string
	TxOptions *txOptions

	Close StepOptions

	Prepare        StepOptions
	PrepareContext StepOptions
	StmtIDKey      string
	StmtOptions    *stmtOptions

	ResetSession StepOptions
	Ping         StepOptions

	ExecContext StepOptions

	QueryContext StepOptions
	RowsOptions  *rowsOptions
}

func defaultConnOptions(driverName string, msgb StepEventMsgBuilder) *connOptions {
	stmtOptions := defaultStmtOptions(msgb)
	rowsOptions := stmtOptions.Rows

	return &connOptions{
		IDGen: IDGeneratorDefault,

		Begin:     *defaultStepOptions(msgb, StepConnBegin, LevelInfo),
		BeginTx:   *defaultStepOptions(msgb, StepConnBeginTx, LevelInfo),
		TxIDKey:   TxIDKeyDefault,
		TxOptions: defaultTxOptions(msgb),

		Close: *defaultStepOptions(msgb, StepConnClose, LevelInfo),

		Prepare:        *defaultStepOptions(msgb, StepConnPrepare, LevelInfo),
		PrepareContext: *defaultStepOptions(msgb, StepConnPrepareContext, LevelInfo),
		StmtIDKey:      StmtIDKeyDefault,
		StmtOptions:    stmtOptions,

		ResetSession: *defaultStepOptions(msgb, StepConnResetSession, LevelTrace),
		Ping:         *defaultStepOptions(msgb, StepConnPing, LevelTrace),

		ExecContext: *defaultStepOptions(msgb, StepConnExecContext, LevelInfo, ConnExecContextErrorHandler(driverName)),

		QueryContext: *defaultStepOptions(msgb, StepConnQueryContext, LevelInfo, ConnQueryContextErrorHandler(driverName)),
		RowsOptions:  rowsOptions,
	}
}

func wrapConn(original driver.Conn, logger *stepLogger, options *connOptions) driver.Conn {
	if original == nil {
		return nil
	}
	switch original.(type) {
	case *connWrapper:
		return original
	// case *connNvcWrapper:
	// 	return original
	case *connWithContextWrapper:
		return original
	case *connNvcWithContextWrapper:
		return original
	}

	connWrapper := connWrapper{original: original, logger: logger, options: options}
	if cwc, ok := original.(connWithContext); ok {
		connWrapper2 := connWithContextWrapper{connWrapper, cwc}
		if nvc, ok := original.(driver.NamedValueChecker); ok {
			return &connNvcWithContextWrapper{connWrapper2, nvc}
		}
		return &connWrapper2
	}

	// Commented out because it's not used.
	//
	// if _, ok := original.(driver.ExecerContext); !ok {
	// 	logger.Warn(fmt.Sprintf("driver.Conn %T does not implement driver.ExecerContext", original))
	// }
	// if _, ok := original.(driver.QueryerContext); !ok {
	// 	logger.Warn(fmt.Sprintf("driver.Conn %T does not implement driver.QueryerContext", original))
	// }
	// if _, ok := original.(driver.ConnPrepareContext); !ok {
	// 	logger.Warn(fmt.Sprintf("driver.Conn %T does not implement driver.ConnPrepareContext", original))
	// }
	// if _, ok := original.(driver.ConnBeginTx); !ok {
	// 	logger.Warn(fmt.Sprintf("driver.Conn %T does not implement driver.ConnBeginTx", original))
	// }

	return &connWrapper
}

// See https://pkg.go.dev/database/sql/driver#pkg-overview

type connWrapper struct {
	original driver.Conn
	logger   *stepLogger
	options  *connOptions
}

// Deprecated interfaces, not implemented.
// var _ driver.Execer = (*conn)(nil)
// var _ driver.Queryer = (*conn)(nil)

var (
	_ driver.Conn      = (*connWrapper)(nil)
	_ driver.Validator = (*connWrapper)(nil)
)

// Begin implements driver.Conn.
func (c *connWrapper) Begin() (driver.Tx, error) {
	var origTx driver.Tx
	attr, err := c.logger.StepWithoutContext(&c.options.Begin, func() (*slog.Attr, error) {
		var err error
		origTx, err = c.original.Begin() //nolint:staticcheck
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(c.options.TxIDKey, c.options.IDGen())
		return &attrRaw, nil
	})
	if err != nil {
		return nil, err
	}
	lg := c.logger
	if attr != nil {
		lg = lg.With(*attr)
	}
	return wrapTx(origTx, lg, c.options.TxOptions), nil
}

// Close implements driver.Conn.
func (c *connWrapper) Close() error {
	return ignoreAttr(c.logger.StepWithoutContext(&c.options.Close, withNilAttr(c.original.Close)))
}

// Prepare implements driver.Conn.
func (c *connWrapper) Prepare(query string) (driver.Stmt, error) {
	var origStmt driver.Stmt
	attr, err := c.logger.With(slog.String("query", query)).StepWithoutContext(&c.options.Prepare, func() (*slog.Attr, error) {
		var err error
		origStmt, err = c.original.Prepare(query)
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(c.options.StmtIDKey, c.options.IDGen())
		return &attrRaw, nil
	})
	if err != nil {
		return nil, err
	}
	lg := c.logger
	if attr != nil {
		lg = lg.With(*attr)
	}
	return wrapStmt(origStmt, lg, c.options.StmtOptions), nil
}

// IsValid implements driver.Validator.
func (c *connWrapper) IsValid() bool {
	// https://cs.opensource.google/go/go/+/master:src/database/sql/sql.go;l=618-621
	if v, ok := c.original.(driver.Validator); ok {
		return v.IsValid()
	}
	return true
}

type connWithContext interface {
	driver.Conn
	driver.ExecerContext
	driver.QueryerContext
	driver.ConnPrepareContext
	driver.ConnBeginTx
}

type connWithContextWrapper struct {
	connWrapper
	originalConn connWithContext
}

// // All Conn implementations should implement the following
// interfaces: Pinger, SessionResetter, and Validator.
var (
	_ driver.Pinger          = (*connWithContextWrapper)(nil)
	_ driver.SessionResetter = (*connWithContextWrapper)(nil)
)

// If named parameters or context are supported, the driver's
// Conn should implement: ExecerContext, QueryerContext,
// ConnPrepareContext, and ConnBeginTx.
var (
	_ driver.ExecerContext      = (*connWithContextWrapper)(nil)
	_ driver.QueryerContext     = (*connWithContextWrapper)(nil)
	_ driver.ConnPrepareContext = (*connWithContextWrapper)(nil)
	_ driver.ConnBeginTx        = (*connWithContextWrapper)(nil)
)

// ResetSession implements driver.SessionResetter.
func (c *connWithContextWrapper) ResetSession(ctx context.Context) error {
	return ignoreAttr(c.logger.Step(ctx, &c.options.ResetSession, func() (*slog.Attr, error) {
		// https://cs.opensource.google/go/go/+/master:src/database/sql/sql.go;l=603-606
		if v, ok := c.original.(driver.SessionResetter); ok {
			return nil, v.ResetSession(ctx)
		}
		return nil, nil
	}))
}

// Ping implements driver.Pinger.
func (c *connWithContextWrapper) Ping(ctx context.Context) error {
	return ignoreAttr(c.logger.Step(ctx, &c.options.Ping, func() (*slog.Attr, error) {
		// https://cs.opensource.google/go/go/+/master:src/database/sql/sql.go;l=882-891
		if p, ok := c.original.(driver.Pinger); ok {
			return nil, p.Ping(ctx)
		}
		return nil, nil
	}))
}

// ExecContext implements driver.ExecerContext.
func (c *connWithContextWrapper) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	var result driver.Result
	lg := c.logger.With(
		slog.String("query", query),
		slog.String("args", fmt.Sprintf("%+v", args)),
	)
	err := ignoreAttr(lg.Step(ctx, &c.options.ExecContext, func() (*slog.Attr, error) {
		var err error
		result, err = c.originalConn.ExecContext(ctx, query, args)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// QueryContext implements driver.QueryerContext.
func (c *connWithContextWrapper) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	var rows driver.Rows
	lg := c.logger.With(
		slog.String("query", query),
		slog.String("args", fmt.Sprintf("%+v", args)),
	)
	err := ignoreAttr(lg.Step(ctx, &c.options.QueryContext, func() (*slog.Attr, error) {
		var err error
		rows, err = c.originalConn.QueryContext(ctx, query, args)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return wrapRows(rows, c.logger, c.options.RowsOptions), nil
}

// PrepareContext implements driver.ConnPrepareContext.
func (c *connWithContextWrapper) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	var stmt driver.Stmt
	attr, err := c.logger.With(slog.String("query", query)).Step(ctx, &c.options.PrepareContext, func() (*slog.Attr, error) {
		var err error
		stmt, err = c.originalConn.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(c.options.StmtIDKey, c.options.IDGen())
		return &attrRaw, nil
	})
	if err != nil {
		return nil, err
	}
	lg := c.logger
	if attr != nil {
		lg = lg.With(*attr)
	}
	return wrapStmt(stmt, lg, c.options.StmtOptions), nil
}

// BeginTx implements driver.ConnBeginTx.
func (c *connWithContextWrapper) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	var tx driver.Tx
	attr, err := c.logger.Step(ctx, &c.options.BeginTx, func() (*slog.Attr, error) {
		var err error
		tx, err = c.originalConn.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(c.options.TxIDKey, c.options.IDGen())
		return &attrRaw, nil
	})
	if err != nil {
		return nil, err
	}
	lg := c.logger
	if attr != nil {
		lg = lg.With(*attr)
	}
	return wrapTx(tx, lg, c.options.TxOptions), nil
}

const driverNameMysql = "mysql"

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

type connNvcWrapper struct {
	connWrapper
	driver.NamedValueChecker
}

var (
	_ driver.Conn              = (*connNvcWrapper)(nil)
	_ driver.NamedValueChecker = (*connNvcWrapper)(nil)
)

type connNvcWithContextWrapper struct {
	connWithContextWrapper
	driver.NamedValueChecker
}

var (
	_ driver.Conn              = (*connNvcWithContextWrapper)(nil)
	_ driver.NamedValueChecker = (*connNvcWithContextWrapper)(nil)
)
