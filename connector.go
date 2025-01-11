package sqlslog

import (
	"context"
	"database/sql/driver"
	"log/slog"
)

type connector struct {
	original driver.Connector
	logger   *slog.Logger
}

var _ driver.Connector = (*connector)(nil)

func wrapConnector(original driver.Connector, logger *slog.Logger) driver.Connector {
	return &connector{original: original, logger: logger}
}

// Connect implements driver.Connector.
func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	lg := c.logger
	lg.Debug("connector.Connect Start")
	origConn, err := c.original.Connect(ctx)
	if err != nil {
		lg.Error("connector.Connect Error", "error", err)
		return nil, err
	}
	lg.Info("connector.Connect Complete")
	return wrapConn(origConn, c.logger), nil
}

// Driver implements driver.Connector.
func (c *connector) Driver() driver.Driver {
	return c.original.Driver()
}
