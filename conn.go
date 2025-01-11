package sqlslog

import (
	"context"
	"database/sql/driver"
	"log/slog"
)

type conn struct {
	original driver.Conn
	logger   *slog.Logger
}

var _ driver.Conn = (*conn)(nil)

// Deprecated interfaces, not implemented.
// var _ driver.Execer = (*conn)(nil)
// var _ driver.Queryer = (*conn)(nil)

// // All Conn implementations should implement the following
// interfaces: Pinger, SessionResetter, and Validator.
var _ driver.Pinger = (*conn)(nil)
var _ driver.SessionResetter = (*conn)(nil)
var _ driver.Validator = (*conn)(nil)

// If named parameters or context are supported, the driver's
// Conn should implement: ExecerContext, QueryerContext,
// ConnPrepareContext, and ConnBeginTx.
var _ driver.ExecerContext = (*conn)(nil)
var _ driver.QueryerContext = (*conn)(nil)
var _ driver.ConnPrepareContext = (*conn)(nil)
var _ driver.ConnBeginTx = (*conn)(nil)

// To support custom data types, implement NamedValueChecker.
// NamedValueChecker also allows queries to accept per-query
// options as a parameter by returning ErrRemoveArgument from CheckNamedValue.
var _ driver.NamedValueChecker = (*conn)(nil)

func wrapConn(original driver.Conn, logger *slog.Logger) *conn {
	return &conn{original: original, logger: logger}
}

// Begin implements driver.Conn.
func (c *conn) Begin() (driver.Tx, error) {
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
func (c *conn) Close() error {
	panic("unimplemented")
}

// Prepare implements driver.Conn.
func (c *conn) Prepare(query string) (driver.Stmt, error) {
	panic("unimplemented")
}

// IsValid implements driver.Validator.
func (c *conn) IsValid() bool {
	panic("unimplemented")
}

// ResetSession implements driver.SessionResetter.
func (c *conn) ResetSession(ctx context.Context) error {
	panic("unimplemented")
}

// Ping implements driver.Pinger.
func (c *conn) Ping(ctx context.Context) error {
	panic("unimplemented")
}

// ExecContext implements driver.ExecerContext.
func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	panic("unimplemented")
}

// QueryContext implements driver.QueryerContext.
func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	panic("unimplemented")
}

// PrepareContext implements driver.ConnPrepareContext.
func (c *conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	panic("unimplemented")
}

// BeginTx implements driver.ConnBeginTx.
func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	panic("unimplemented")
}

// CheckNamedValue implements driver.NamedValueChecker.
func (c *conn) CheckNamedValue(*driver.NamedValue) error {
	panic("unimplemented")
}
