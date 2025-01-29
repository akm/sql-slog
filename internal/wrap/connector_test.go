package wrap

import (
	"context"
	"database/sql/driver"
	"testing"
)

type mockConnectorForWrapConnector struct{}

var _ driver.Connector = (*mockConnectorForWrapConnector)(nil)

func (m *mockConnectorForWrapConnector) Connect(context.Context) (driver.Conn, error) {
	panic("unimplemented")
}

func (m *mockConnectorForWrapConnector) Driver() driver.Driver {
	return nil
}

func TestConnectorDriver(t *testing.T) {
	t.Parallel()
	mock := &mockConnectorForWrapConnector{}
	logger := &SQLLogger{}
	conn := WrapConnector(mock, logger)
	if conn.Driver() != nil {
		t.Fatal("Expected nil")
	}
}
