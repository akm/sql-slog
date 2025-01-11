package sqlslog

import (
	"context"
	"database/sql/driver"
	"log/slog"
)

type connWrapper struct {
	original driver.Conn
	logger   *slog.Logger
}

var _ driver.Conn = (*connWrapper)(nil)

// Deprecated interfaces, not implemented.
// var _ driver.Execer = (*conn)(nil)
// var _ driver.Queryer = (*conn)(nil)

// // All Conn implementations should implement the following
// interfaces: Pinger, SessionResetter, and Validator.
var _ driver.Pinger = (*connWrapper)(nil)
var _ driver.SessionResetter = (*connWrapper)(nil)
var _ driver.Validator = (*connWrapper)(nil)

// If named parameters or context are supported, the driver's
// Conn should implement: ExecerContext, QueryerContext,
// ConnPrepareContext, and ConnBeginTx.
var _ driver.ExecerContext = (*connWrapper)(nil)
var _ driver.QueryerContext = (*connWrapper)(nil)
var _ driver.ConnPrepareContext = (*connWrapper)(nil)
var _ driver.ConnBeginTx = (*connWrapper)(nil)

// To support custom data types, implement NamedValueChecker.
// NamedValueChecker also allows queries to accept per-query
// options as a parameter by returning ErrRemoveArgument from CheckNamedValue.
var _ driver.NamedValueChecker = (*connWrapper)(nil)

func wrapConn(original driver.Conn, logger *slog.Logger) *connWrapper {
	return &connWrapper{original: original, logger: logger}
}

// Begin implements driver.Conn.
func (c *connWrapper) Begin() (driver.Tx, error) {
	var origTx driver.Tx
	err := logAction(c.logger, "Begin", func() error {
		var err error
		origTx, err = c.original.Begin()
		return err
	})
	if err != nil {
		return nil, err
	}
	return wrapTx(origTx, c.logger), nil
}

// Close implements driver.Conn.
func (c *connWrapper) Close() error {
	return logAction(c.logger, "Close", c.original.Close)
}

// Prepare implements driver.Conn.
func (c *connWrapper) Prepare(query string) (driver.Stmt, error) {
	lg := c.logger.With(slog.String("query", query))
	var origStmt driver.Stmt
	err := logAction(lg, "Prepare", func() error {
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

// ResetSession implements driver.SessionResetter.
func (c *connWrapper) ResetSession(ctx context.Context) error {
	return logActionContext(ctx, c.logger, "ResetSession", func() error {
		// https://cs.opensource.google/go/go/+/master:src/database/sql/sql.go;l=603-606
		if v, ok := c.original.(driver.SessionResetter); ok {
			return v.ResetSession(ctx)
		}
		return nil
	})
}

// Ping implements driver.Pinger.
func (c *connWrapper) Ping(ctx context.Context) error {
	return logActionContext(ctx, c.logger, "Ping", func() error {
		// https://cs.opensource.google/go/go/+/master:src/database/sql/sql.go;l=882-891
		if p, ok := c.original.(driver.Pinger); ok {
			return p.Ping(ctx)
		}
		return nil
	})
}

// ExecContext implements driver.ExecerContext.
func (c *connWrapper) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	panic("unimplemented")
}

// QueryContext implements driver.QueryerContext.
func (c *connWrapper) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	panic("unimplemented")
}

// PrepareContext implements driver.ConnPrepareContext.
func (c *connWrapper) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	panic("unimplemented")
}

// BeginTx implements driver.ConnBeginTx.
func (c *connWrapper) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	panic("unimplemented")
}

// CheckNamedValue implements driver.NamedValueChecker.
func (c *connWrapper) CheckNamedValue(*driver.NamedValue) error {
	panic("unimplemented")
}
