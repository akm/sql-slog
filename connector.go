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
	err := c.logger.Step(ctx, &c.logger.options.connectorConnect, func() error {
		var err error
		origConn, err = c.original.Connect(ctx)
		return err
	})
	if err != nil {
		return nil, err
	}
	return wrapConn(origConn, c.logger), nil
}

// Driver implements driver.Connector.
func (c *connector) Driver() driver.Driver {
	return c.original.Driver()
}

// ConnectorConnectErrorHandler returns a function that handles the error of driver.Connector.Connect.
// The function returns completed and slice of slog.Attr.
//
// # For Postgres:
// If err is nil, it returns true and a slice of slog.Attr{slog.Bool("eof", false)}.
// If err is io.EOF, it returns true and a slice of slog.Attr{slog.Bool("eof", true)}.
// Otherwise, it returns false and nil.
func ConnectorConnectErrorHandler(driverName string) func(err error) (bool, []slog.Attr) {
	switch driverName {
	case "postgres":
		return func(err error) (bool, []slog.Attr) {
			if err == nil {
				return true, []slog.Attr{slog.Bool("eof", false)}
			}
			if errors.Is(err, io.EOF) {
				return true, []slog.Attr{slog.Bool("eof", true)}
			}
			return false, nil
		}
	default:
		return nil
	}
}
