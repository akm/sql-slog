package sqlslog_test

import (
	"database/sql"
	"database/sql/driver"
)

type MockDriver struct {
	OpenResult driver.Conn
	OpenError  error
}

func (m *MockDriver) Open(string) (driver.Conn, error) {
	return m.OpenResult, m.OpenError
}

var mockDriver driver.Driver = &MockDriver{}

type MockConn struct{}

func (m *MockConn) Begin() (driver.Tx, error) {
	return nil, nil // nolint: nilnil
}

func (m *MockConn) Close() error {
	return nil
}

func (m *MockConn) Prepare(string) (driver.Stmt, error) {
	return nil, nil // nolint: nilnil
}

var _ driver.Conn = (*MockConn)(nil)

func init() { // nolint: gochecknoinits
	sql.Register("mock", mockDriver)
}
