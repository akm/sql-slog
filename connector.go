package sqlslog

import (
	"context"
	"database/sql/driver"
	"errors"
	"io"
	"log/slog"
)

type connector struct {
	original driver.Connector
	logger   *logger
}

var _ driver.Connector = (*connector)(nil)

func wrapConnector(original driver.Connector, logger *logger) driver.Connector {
	return &connector{original: original, logger: logger}
}

// Connect implements driver.Connector.
func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	var origConn driver.Conn
	err := ignoreAttr(c.logger.Step(ctx, &c.logger.options.connectorConnect, func() (*slog.Attr, error) {
		var err error
		origConn, err = c.original.Connect(ctx)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	opts := c.logger.options
	return wrapConn(origConn, c.logger, &connOptions{
		idGen:   opts.idGen,
		Begin:   &opts.connBegin,
		BeginTx: &opts.connBeginTx,
		txIDKey: opts.txIDKey,
		Tx: &txOptions{
			Commit:   &opts.txCommit,
			Rollback: &opts.txRollback,
		},
		Close:          &opts.connClose,
		Prepare:        &opts.connPrepare,
		PrepareContext: &opts.connPrepareContext,
		stmtIDKey:      opts.stmtIDKey,
		Stmt: &stmtOptions{
			Close:        &opts.stmtClose,
			Exec:         &opts.stmtExec,
			Query:        &opts.stmtQuery,
			ExecContext:  &opts.stmtExecContext,
			QueryContext: &opts.stmtQueryContext,
			Rows: &rowsOptions{
				Close:         &opts.rowsClose,
				Next:          &opts.rowsNext,
				NextResultSet: &opts.rowsNextResultSet,
			},
		},
		ResetSession: &opts.connResetSession,
		Ping:         &opts.connPing,
		ExecContext:  &opts.connExecContext,
		QueryContext: &opts.connQueryContext,
		Rows: &rowsOptions{
			Close:         &opts.rowsClose,
			Next:          &opts.rowsNext,
			NextResultSet: &opts.rowsNextResultSet,
		},
	}), nil
}

// Driver implements driver.Connector.
func (c *connector) Driver() driver.Driver {
	return c.original.Driver()
}

// ConnectorConnectErrorHandler returns a function that handles errors from driver.Connector.Connect.
// The function returns a boolean indicating completion and a slice of slog.Attr.
//
// # For Postgres:
// If err is nil, it returns true and a slice of slog.Attr{slog.Bool("success", true)}.
// If err is io.EOF, it returns true and a slice of slog.Attr{slog.Bool("success", false)}.
// Otherwise, it returns false and nil.
func ConnectorConnectErrorHandler(driverName string) func(err error) (bool, []slog.Attr) {
	switch driverName {
	case "mysql":
		return func(err error) (bool, []slog.Attr) {
			if err == nil {
				return true, []slog.Attr{slog.Bool("success", true)}
			}
			if err.Error() == "driver: bad connection" {
				return true, []slog.Attr{slog.Bool("success", false)}
			}
			return false, nil
		}
	case "postgres":
		return func(err error) (bool, []slog.Attr) {
			if err == nil {
				return true, []slog.Attr{slog.Bool("success", true)}
			}
			if errors.Is(err, io.EOF) {
				return true, []slog.Attr{slog.Bool("success", false)}
			}
			return false, nil
		}
	default:
		return nil
	}
}
