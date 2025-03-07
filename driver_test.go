package sqlslog

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"io"
	"log/slog"
	"testing"
)

type mockErrorDiverContext struct{}

// Open implements driver.Driver.
func (m *mockErrorDiverContext) Open(string) (driver.Conn, error) {
	return nil, errors.New("unexpected error")
}

// OpenConnector implements driver.DriverContext.
func (m *mockErrorDiverContext) OpenConnector(string) (driver.Connector, error) {
	return nil, errors.New("unexpected error")
}

var (
	_ driver.Driver        = (*mockErrorDiverContext)(nil)
	_ driver.DriverContext = (*mockErrorDiverContext)(nil)
)

func TestDriverContextWrapperOpenConnector(t *testing.T) {
	t.Parallel()
	t.Run("unexpected error", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		logger := slog.New(NewTextHandler(buf, nil))
		dw := wrapDriver(&mockErrorDiverContext{},
			newStepLogger(logger, defaultStepLoggerOptions()),
			defaultDriverOptions("sqlite3", StepEventMsgWithoutEventName),
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
	t.Parallel()
	t.Run("postgres", func(t *testing.T) {
		t.Parallel()
		t.Run("no error", func(t *testing.T) {
			t.Parallel()
			eh := DriverOpenErrorHandler("postgres")
			completed, attrs := eh(nil)
			if !completed {
				t.Error("expected completed to be false")
			}
			if attrs == nil {
				t.Error("expected attrs not to be nil")
			}
		})
		t.Run("unexpected error", func(t *testing.T) {
			t.Parallel()
			eh := DriverOpenErrorHandler("postgres")
			completed, attrs := eh(errors.New("unexpected error"))
			if completed {
				t.Error("expected completed to be false")
			}
			if attrs != nil {
				t.Error("expected attrs to be nil")
			}
		})
		t.Run("io.EOF", func(t *testing.T) {
			t.Parallel()
			eh := DriverOpenErrorHandler("postgres")
			completed, attrs := eh(io.EOF)
			if !completed {
				t.Errorf("expected completed to be true")
			}
			if attrs == nil {
				t.Error("expected attrs to be non-nil")
			}
		})
	})
	t.Run("mysql", func(t *testing.T) {
		t.Parallel()
		eh := DriverOpenErrorHandler("mysql")
		if eh != nil {
			t.Error("expected to be nil")
		}
	})
}
