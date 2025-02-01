package sqlslog

import (
	"context"
	"database/sql/driver"
	"errors"
	"io"
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

	t.Run("postgres io.EOF", func(t *testing.T) {
		t.Parallel()
		errHandler := ConnectorConnectErrorHandler("postgres")
		complete, attrs := errHandler(io.EOF)
		if !complete {
			t.Fatal("Expected true")
		}
		if attrs == nil {
			t.Fatal("Expected non-nil")
		}
	})
}

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
	logger := &logger{}
	conn := WrapConnector(mock, logger, nil)
	if conn.Driver() != nil {
		t.Fatal("Expected nil")
	}
}
