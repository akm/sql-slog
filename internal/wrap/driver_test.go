package wrap

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"log/slog"
	"testing"

	"github.com/akm/sql-slog/internal/opts"
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

		driverOptions := opts.DefaultDriverOptions("sqlite3", opts.StepLogMsgWithoutEventName)
		buf := bytes.NewBuffer(nil)
		logger := slog.New(NewTextHandler(buf, nil))
		dw := WrapDriver(&mockErrorDiverContext{},
			NewStepLogger(logger, DurationAttrFunc(opts.DurationKeyDefault, opts.DurationNanoSeconds)),
			driverOptions,
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
