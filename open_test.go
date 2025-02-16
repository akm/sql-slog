package sqlslog

import (
	"context"
	"database/sql/driver"
	"errors"
	"log/slog"
	"testing"
)

func TestOpen(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()

	db, _, err := Open(ctx, "invalid-driver", "")
	if err == nil {
		t.Fatal("Expected error")
	}
	if err.Error() != "sql: unknown driver \"invalid-driver\" (forgotten import?)" {
		t.Fatalf("Unexpected error: %v", err)
	}
	if db != nil {
		t.Fatal("Expected nil db")
	}
}

type errorDriverContext struct {
	error error
}

var (
	_ driver.Driver        = (*errorDriverContext)(nil)
	_ driver.DriverContext = (*errorDriverContext)(nil)
)

func newErrorDriverContext(err error) *errorDriverContext {
	return &errorDriverContext{error: err}
}

// OpenConnector implements driver.DriverContext.
func (e *errorDriverContext) OpenConnector(string) (driver.Connector, error) {
	return nil, e.error
}

// Open implements driver.Driver.
func (e *errorDriverContext) Open(string) (driver.Conn, error) {
	return nil, e.error
}

func TestOpenWithDriver(t *testing.T) {
	t.Parallel()
	t.Run("unknown error", func(t *testing.T) {
		t.Parallel()
		t.Run("DriverContext", func(t *testing.T) {
			drv := newErrorDriverContext(errors.New("unknown error"))
			stepLogger := newStepLogger(newStepLoggerOptions(slog.New(slog.NewJSONHandler(nil, nil))))
			if _, err := openWithWrappedDriver(drv, "invalid-dsn", stepLogger, nil); err == nil {
				t.Fatal("Expected error")
			}
		})
	})
}
