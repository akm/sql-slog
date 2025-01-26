package sqlslog

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"
)

func TestConnectorConnectErrorHandler(t *testing.T) {
	t.Parallel()
	testcases := []string{
		"mysql",
		"postgres",
	}
	for _, driverName := range testcases {
		t.Run(driverName, func(t *testing.T) {
			t.Parallel()
			errHandler := ConnectorConnectErrorHandler(driverName)
			complete, attrs := errHandler(errors.New("dummy"))
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
	t.Parallel()
	mock := &mockConnectorForWrapConnector{}
	logger := &logger{}
	conn := wrapConnector(mock, logger)
	if conn.Driver() != nil {
		t.Fatal("Expected nil")
	}
}
