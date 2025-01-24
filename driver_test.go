package sqlslog

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"log/slog"
	"testing"
)

type mockErrorDiverContext struct {
}

// Open implements driver.Driver.
func (m *mockErrorDiverContext) Open(name string) (driver.Conn, error) {
	return nil, fmt.Errorf("unexpected error")
}

// OpenConnector implements driver.DriverContext.
func (m *mockErrorDiverContext) OpenConnector(name string) (driver.Connector, error) {
	return nil, fmt.Errorf("unexpected error")
}

var _ driver.Driver = (*mockErrorDiverContext)(nil)
var _ driver.DriverContext = (*mockErrorDiverContext)(nil)

func TestDriverContextWrapperOpenConnector(t *testing.T) {
	t.Run("unexpected error", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		logger := slog.New(NewTextHandler(buf, nil))
		dw := wrapDriver(&mockErrorDiverContext{},
			newLogger(logger, newOptions("sqlite3")),
		)
		dwc, ok := dw.(driver.DriverContext)
		if !ok {
			t.Fatal("expected to be driver.DriverContext")
		}
		_, err := dwc.OpenConnector("dsn")
		if err == nil {
			t.Fatal("expected error to be not nil")
		}
	})
}

func TestDriverOpenErrorHandler(t *testing.T) {
	t.Run("postgres", func(t *testing.T) {
		t.Run("unexpected error", func(t *testing.T) {
			eh := DriverOpenErrorHandler("postgres")
			completed, attrs := eh(errors.New("unexpected error"))
			if completed {
				t.Error("expected completed to be false")
			}
			if attrs != nil {
				t.Error("expected attrs to be nil")
			}
		})
	})

}
