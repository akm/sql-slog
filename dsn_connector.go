package sqlslog

import (
	"context"
	"database/sql/driver"
)

type dsnConnector struct {
	dsn    string
	driver driver.Driver
}

var _ driver.Connector = (*dsnConnector)(nil)

func (t dsnConnector) Connect(_ context.Context) (driver.Conn, error) {
	return t.driver.Open(t.dsn)
}

func (t dsnConnector) Driver() driver.Driver {
	return t.driver
}
