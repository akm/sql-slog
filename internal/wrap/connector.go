package wrap

import (
	"context"
	"database/sql/driver"
	"log/slog"
)

type connector struct {
	original driver.Connector
	logger   *SqlLogger
}

var _ driver.Connector = (*connector)(nil)

func WrapConnector(original driver.Connector, logger *SqlLogger) driver.Connector {
	return &connector{original: original, logger: logger}
}

// Connect implements driver.Connector.
func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	var origConn driver.Conn
	err := IgnoreAttr(c.logger.Step(ctx, &c.logger.Options.ConnectorConnect, func() (*slog.Attr, error) {
		var err error
		origConn, err = c.original.Connect(ctx)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return WrapConn(origConn, c.logger), nil
}

// Driver implements driver.Connector.
func (c *connector) Driver() driver.Driver {
	return c.original.Driver()
}
