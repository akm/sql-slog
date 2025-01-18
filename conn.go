package sqlslog

import (
	"context"
	"database/sql/driver"
	"fmt"
	"log/slog"
)

func wrapConn(original driver.Conn, logger *logger) driver.Conn {
	if original == nil {
		return nil
	}
	switch original.(type) {
	case *connWithContextWrapper:
		return original
		// Commented out because it's not used.
		// case *connWrapper:
		// 	return original
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
	logger   *logger
}

// Deprecated interfaces, not implemented.
// var _ driver.Execer = (*conn)(nil)
// var _ driver.Queryer = (*conn)(nil)

var _ driver.Conn = (*connWrapper)(nil)
var _ driver.Validator = (*connWrapper)(nil)

// To support custom data types, implement NamedValueChecker.
// NamedValueChecker also allows queries to accept per-query
// options as a parameter by returning ErrRemoveArgument from CheckNamedValue.
var _ driver.NamedValueChecker = (*connWrapper)(nil)

// Begin implements driver.Conn.
func (c *connWrapper) Begin() (driver.Tx, error) {
	var origTx driver.Tx
	err := c.logger.logAction(&c.logger.options.connBegin, func() error {
		var err error
		origTx, err = c.original.Begin() //nolint:staticcheck
		return err
	})
	if err != nil {
		return nil, err
	}
	return wrapTx(origTx, c.logger), nil
}

// Close implements driver.Conn.
func (c *connWrapper) Close() error {
	return c.logger.logAction(&c.logger.options.connClose, c.original.Close)
}

// Prepare implements driver.Conn.
func (c *connWrapper) Prepare(query string) (driver.Stmt, error) {
	lg := c.logger.With(slog.String("query", query))
	var origStmt driver.Stmt
	err := lg.logAction(&c.logger.options.connPrepare, func() error {
		var err error
		origStmt, err = c.original.Prepare(query)
		return err
	})
	if err != nil {
		return nil, err
	}
	return wrapStmt(origStmt, lg), nil
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
var _ driver.Pinger = (*connWithContextWrapper)(nil)
var _ driver.SessionResetter = (*connWithContextWrapper)(nil)

// If named parameters or context are supported, the driver's
// Conn should implement: ExecerContext, QueryerContext,
// ConnPrepareContext, and ConnBeginTx.
var _ driver.ExecerContext = (*connWithContextWrapper)(nil)
var _ driver.QueryerContext = (*connWithContextWrapper)(nil)
var _ driver.ConnPrepareContext = (*connWithContextWrapper)(nil)
var _ driver.ConnBeginTx = (*connWithContextWrapper)(nil)

// ResetSession implements driver.SessionResetter.
func (c *connWithContextWrapper) ResetSession(ctx context.Context) error {
	return c.logger.logActionContext(ctx, &c.logger.options.connResetSession, func() error {
		// https://cs.opensource.google/go/go/+/master:src/database/sql/sql.go;l=603-606
		if v, ok := c.original.(driver.SessionResetter); ok {
			return v.ResetSession(ctx)
		}
		return nil
	})
}

// Ping implements driver.Pinger.
func (c *connWithContextWrapper) Ping(ctx context.Context) error {
	return c.logger.logActionContext(ctx, &c.logger.options.connPing, func() error {
		// https://cs.opensource.google/go/go/+/master:src/database/sql/sql.go;l=882-891
		if p, ok := c.original.(driver.Pinger); ok {
			return p.Ping(ctx)
		}
		return nil
	})
}

// ExecContext implements driver.ExecerContext.
func (c *connWithContextWrapper) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	var result driver.Result
	lg := c.logger.With(
		slog.String("query", query),
		slog.String("args", fmt.Sprintf("%+v", args)),
	)
	err := lg.logActionContext(ctx, &c.logger.options.connExecContext, func() error {
		var err error
		result, err = c.originalConn.ExecContext(ctx, query, args)
		return err
	})
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
	err := lg.logActionContext(ctx, &c.logger.options.connQueryContext, func() error {
		var err error
		rows, err = c.originalConn.QueryContext(ctx, query, args)
		return err
	})
	if err != nil {
		return nil, err
	}
	return wrapRows(rows, c.logger), nil
}

// PrepareContext implements driver.ConnPrepareContext.
func (c *connWithContextWrapper) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	var stmt driver.Stmt
	lg := c.logger.With(slog.String("query", query))
	err := lg.logActionContext(ctx, &c.logger.options.connPrepareContext, func() error {
		var err error
		stmt, err = c.originalConn.PrepareContext(ctx, query)
		return err
	})
	if err != nil {
		return nil, err
	}
	return wrapStmt(stmt, c.logger), nil
}

// BeginTx implements driver.ConnBeginTx.
func (c *connWithContextWrapper) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	var tx driver.Tx
	err := c.logger.logActionContext(ctx, &c.logger.options.connBeginTx, func() error {
		var err error
		tx, err = c.originalConn.BeginTx(ctx, opts)
		return err
	})
	if err != nil {
		return nil, err
	}
	return wrapTx(tx, c.logger), nil
}
