package sqlslog

import (
	"context"
	"database/sql/driver"
	"testing"
)

type mockDriverForDsnConnector struct {
	resultConn driver.Conn
	resultErr  error
}

var _ driver.Driver = (*mockDriverForDsnConnector)(nil)

func (m *mockDriverForDsnConnector) Open(string) (driver.Conn, error) {
	return m.resultConn, m.resultErr
}

func TestDsnConnector(t *testing.T) {
	t.Parallel()
	var d dsnConnector
	dsn := "dsn"
	drv := &mockDriverForDsnConnector{}
	d = dsnConnector{dsn: dsn, driver: drv}

	t.Run("Connect", func(t *testing.T) {
		t.Parallel()
		_, _ = d.Connect(context.Background())
	})

	t.Run("Driver", func(t *testing.T) {
		t.Parallel()
		if d.Driver() != drv {
			t.Errorf("expected %v, got %v", drv, d.Driver())
		}
	})
}
