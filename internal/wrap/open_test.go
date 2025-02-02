package wrap

import (
	"bytes"
	"context"
	"database/sql/driver"
	"errors"
	"log/slog"
	"testing"
)

func TestOpen(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	logger := slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	db, err := Open(ctx, "invalid-driver", "", Logger(logger))
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
			stepLogger := NewStepLogger(slog.New(slog.NewJSONHandler(nil, nil)), nil)
			_, err := openWithDriver(drv, "invalid-dsn", stepLogger, nil)
			if err == nil {
				t.Fatal("Expected error")
			}
		})
	})
}
