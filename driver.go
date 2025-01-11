package sqlslog

import (
	"context"
	"database/sql/driver"
)

// https://pkg.go.dev/database/sql/driver@go1.23.4#pkg-overview
// The driver interface has evolved over time. Drivers
// should implement Connector and DriverContext interfaces.
type driverImpl struct {
}

var _ driver.Driver = (*driverImpl)(nil)
var _ driver.DriverContext = (*driverImpl)(nil)
var _ driver.Connector = (*driverImpl)(nil)

// Open implements driver.Driver.
func (c *driverImpl) Open(name string) (driver.Conn, error) {
	panic("unimplemented")
}

// OpenConnector implements driver.DriverContext.
func (c *driverImpl) OpenConnector(name string) (driver.Connector, error) {
	panic("unimplemented")
}

// Connect implements driver.Connector.
func (c *driverImpl) Connect(context.Context) (driver.Conn, error) {
	panic("unimplemented")
}

// Driver implements driver.Connector.
func (c *driverImpl) Driver() driver.Driver {
	panic("unimplemented")
}
