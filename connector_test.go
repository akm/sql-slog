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
		errHandler := ConnectorConnectErrorHandler(driverName)
		t.Run(driverName, func(t *testing.T) {
			t.Parallel()
			t.Run("no error", func(t *testing.T) {
				complete, attrs := errHandler(nil)
				if !complete {
					t.Error("Expected true")
				}
				if len(attrs) < 1 {
					t.Error("Expected non-empty")
				}
			})
			t.Run("unexpected error", func(t *testing.T) {
				t.Parallel()
				complete, attrs := errHandler(errors.New("unexpected-error"))
				if complete {
					t.Fatal("Expected false")
				}
				if attrs != nil {
					t.Fatal("Expected nil")
				}
			})
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
	t.Run("mysql driver: bad connection", func(t *testing.T) {
		t.Parallel()
		errHandler := ConnectorConnectErrorHandler("mysql")
		complete, attrs := errHandler(errors.New("driver: bad connection"))
		if !complete {
			t.Fatal("Expected true")
		}
		if attrs == nil {
			t.Fatal("Expected non-nil")
		}
	})
	t.Run("sqlite3", func(t *testing.T) {
		t.Parallel()
		errHandler := ConnectorConnectErrorHandler("sqlite3")
		if errHandler != nil {
			t.Fatal("Expected nil")
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
	logger := &stepLogger{}
	connectorOptions := defaultConnectorOptions("dummy", StepEventMsgWithoutEventName)
	conn := wrapConnector(mock, logger, connectorOptions)
	if conn.Driver() != nil {
		t.Fatal("Expected nil")
	}
}
