package wrap

import (
	"context"
	"database/sql/driver"
	"fmt"
	"log/slog"
)

func WrapConn(original driver.Conn, logger *SqlLogger) driver.Conn {
	if original == nil {
		return nil
	}
	if _, ok := original.(*connWithContextWrapper); ok {
		return original
	}

	if cwc, ok := original.(connWithContext); ok {
		return &connWithContextWrapper{connWrapper{original: original, logger: logger}, cwc}
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

	return &connWrapper{original: original, logger: logger}
}

// See https://pkg.go.dev/database/sql/driver#pkg-overview

type connWrapper struct {
	original driver.Conn
	logger   *SqlLogger
}

// Deprecated interfaces, not implemented.
// var _ driver.Execer = (*conn)(nil)
// var _ driver.Queryer = (*conn)(nil)

var (
	_ driver.Conn      = (*connWrapper)(nil)
	_ driver.Validator = (*connWrapper)(nil)
)

// To support custom data types, implement NamedValueChecker.
// NamedValueChecker also allows queries to accept per-query
// options as a parameter by returning ErrRemoveArgument from CheckNamedValue.
var _ driver.NamedValueChecker = (*connWrapper)(nil)

// Begin implements driver.Conn.
func (c *connWrapper) Begin() (driver.Tx, error) {
	var origTx driver.Tx
	attr, err := c.logger.StepWithoutContext(&c.logger.Options.ConnBegin, func() (*slog.Attr, error) {
		var err error
		origTx, err = c.original.Begin() //nolint:staticcheck
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(c.logger.Options.TxIDKey, c.logger.Options.IDGen())
		return &attrRaw, nil
	})
	if err != nil {
		return nil, err
	}
	lg := c.logger
	if attr != nil {
		lg = lg.With(*attr)
	}
	return WrapTx(origTx, lg), nil
}

// Close implements driver.Conn.
func (c *connWrapper) Close() error {
	return IgnoreAttr(c.logger.StepWithoutContext(&c.logger.Options.ConnClose, WithNilAttr(c.original.Close)))
}

// Prepare implements driver.Conn.
func (c *connWrapper) Prepare(query string) (driver.Stmt, error) {
	var origStmt driver.Stmt
	attr, err := c.logger.With(slog.String("query", query)).StepWithoutContext(&c.logger.Options.ConnPrepare, func() (*slog.Attr, error) {
		var err error
		origStmt, err = c.original.Prepare(query)
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(c.logger.Options.StmtIDKey, c.logger.Options.IDGen())
		return &attrRaw, nil
	})
	if err != nil {
		return nil, err
	}
	lg := c.logger
	if attr != nil {
		lg = lg.With(*attr)
	}
	return WrapStmt(origStmt, lg), nil
}

// IsValid implements driver.Validator.
func (c *connWrapper) IsValid() bool {
	// https://cs.opensource.google/go/go/+/master:src/database/sql/sql.go;l=618-621
	if v, ok := c.original.(driver.Validator); ok {
		return v.IsValid()
	}
	return true
}

// CheckNamedValue implements driver.NamedValueChecker.
func (c *connWrapper) CheckNamedValue(namedValue *driver.NamedValue) error {
	if v, ok := c.original.(driver.NamedValueChecker); ok {
		return v.CheckNamedValue(namedValue)
	}
	return nil
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
	return IgnoreAttr(c.logger.Step(ctx, &c.logger.Options.ConnResetSession, func() (*slog.Attr, error) {
		// https://cs.opensource.google/go/go/+/master:src/database/sql/sql.go;l=603-606
		if v, ok := c.original.(driver.SessionResetter); ok {
			return nil, v.ResetSession(ctx)
		}
		return nil, nil
	}))
}

// Ping implements driver.Pinger.
func (c *connWithContextWrapper) Ping(ctx context.Context) error {
	return IgnoreAttr(c.logger.Step(ctx, &c.logger.Options.ConnPing, func() (*slog.Attr, error) {
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
	err := IgnoreAttr(lg.Step(ctx, &c.logger.Options.ConnExecContext, func() (*slog.Attr, error) {
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
	err := IgnoreAttr(lg.Step(ctx, &c.logger.Options.ConnQueryContext, func() (*slog.Attr, error) {
		var err error
		rows, err = c.originalConn.QueryContext(ctx, query, args)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return WrapRows(rows, c.logger), nil
}

// PrepareContext implements driver.ConnPrepareContext.
func (c *connWithContextWrapper) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	var stmt driver.Stmt
	attr, err := c.logger.With(slog.String("query", query)).Step(ctx, &c.logger.Options.ConnPrepareContext, func() (*slog.Attr, error) {
		var err error
		stmt, err = c.originalConn.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(c.logger.Options.StmtIDKey, c.logger.Options.IDGen())
		return &attrRaw, nil
	})
	if err != nil {
		return nil, err
	}
	lg := c.logger
	if attr != nil {
		lg = lg.With(*attr)
	}
	return WrapStmt(stmt, lg), nil
}

// BeginTx implements driver.ConnBeginTx.
func (c *connWithContextWrapper) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	var tx driver.Tx
	attr, err := c.logger.Step(ctx, &c.logger.Options.ConnBeginTx, func() (*slog.Attr, error) {
		var err error
		tx, err = c.originalConn.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(c.logger.Options.TxIDKey, c.logger.Options.IDGen())
		return &attrRaw, nil
	})
	if err != nil {
		return nil, err
	}
	lg := c.logger
	if attr != nil {
		lg = lg.With(*attr)
	}
	return WrapTx(tx, lg), nil
}
