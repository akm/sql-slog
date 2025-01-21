package sqlslog

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"
)

func TestConnectorConnectErrorHandler(t *testing.T) {
	testcases := []string{
		"mysql",
		"postgres",
	}
	for _, driverName := range testcases {
		t.Run(driverName, func(t *testing.T) {
			errHandler := ConnectorConnectErrorHandler(driverName)
			complete, attrs := errHandler(fmt.Errorf("dummy"))
			if complete {
				t.Fatal("Expected false")
			}
			if attrs != nil {
				t.Fatal("Expected nil")
			}
		})
	}
}

type mockConnectorForWrapConnector struct {
}

var _ driver.Connector = (*mockConnectorForWrapConnector)(nil)

func (m *mockConnectorForWrapConnector) Connect(context.Context) (driver.Conn, error) {
	panic("unimplemented")
}

func (m *mockConnectorForWrapConnector) Driver() driver.Driver {
	return nil
}

func TestConnectorDriver(t *testing.T) {
	mock := &mockConnectorForWrapConnector{}
	logger := &logger{}
	conn := wrapConnector(mock, logger)
	if conn.Driver() != nil {
		t.Fatal("Expected nil")
	}
}
