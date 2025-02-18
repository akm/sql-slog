package sqlslog_test

import (
	"database/sql"
	"database/sql/driver"
)

type MockDriver struct {
	OpenResult driver.Conn
	OpenError  error
}

func (m *MockDriver) Open(name string) (driver.Conn, error) {
	return m.OpenResult, m.OpenError
}

var (
	mockDriver driver.Driver = &MockDriver{}
)

type MockConn struct {
}

func (m *MockConn) Begin() (driver.Tx, error) {
	return nil, nil
}

func (m *MockConn) Close() error {
	return nil
}

func (m *MockConn) Prepare(query string) (driver.Stmt, error) {
	return nil, nil
}

var _ driver.Conn = (*MockConn)(nil)

func init() {
	sql.Register("mock", mockDriver)
}
