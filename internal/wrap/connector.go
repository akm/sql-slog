package wrap

import (
	"context"
	"database/sql/driver"
	"log/slog"

	"github.com/akm/sql-slog/opts"
)

type ConnectorOptions = opts.ConnectorOptions

type connector struct {
	original driver.Connector
	logger   *StepLogger
	options  *ConnectorOptions
}

var _ driver.Connector = (*connector)(nil)

func WrapConnector(original driver.Connector, logger *StepLogger, options *ConnectorOptions) driver.Connector { //nolint:revive
	return &connector{original: original, logger: logger, options: options}
}

// Connect implements driver.Connector.
func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	var origConn driver.Conn
	err := IgnoreAttr(c.logger.Step(ctx, c.options.Connect, func() (*slog.Attr, error) {
		var err error
		origConn, err = c.original.Connect(ctx)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return WrapConn(origConn, c.logger, c.options.Conn), nil
}

// Driver implements driver.Connector.
func (c *connector) Driver() driver.Driver {
	return c.original.Driver()
}
