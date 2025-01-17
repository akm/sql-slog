package sqlslog

import (
	"context"
	"database/sql/driver"
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
	err := c.logger.logActionContext(ctx, "Connector.Connect", func() error {
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
