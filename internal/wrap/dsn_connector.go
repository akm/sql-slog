package wrap

import (
	"context"
	"database/sql/driver"
)

type DsnConnector struct {
	dsn    string
	driver driver.Driver
}

var _ driver.Connector = (*DsnConnector)(nil)

func NewDsnConnector(dsn string, driver driver.Driver) *DsnConnector {
	return &DsnConnector{dsn: dsn, driver: driver}
}

// Connect implements driver.Connector.
func (t DsnConnector) Connect(_ context.Context) (driver.Conn, error) {
	return t.driver.Open(t.dsn)
}

// Driver implements driver.Connector.
func (t DsnConnector) Driver() driver.Driver {
	return t.driver
}
