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
		opts := newOptions("sqlite3")
		connOptions := &connOptions{
			idGen:   opts.idGen,
			Begin:   &opts.connBegin,
			BeginTx: &opts.connBeginTx,
			txIDKey: opts.txIDKey,
			Tx: &txOptions{
				Commit:   &opts.txCommit,
				Rollback: &opts.txRollback,
			},
			Close:          &opts.connClose,
			Prepare:        &opts.connPrepare,
			PrepareContext: &opts.connPrepareContext,
			stmtIDKey:      opts.stmtIDKey,
			Stmt: &stmtOptions{
				Close:        &opts.stmtClose,
				Exec:         &opts.stmtExec,
				Query:        &opts.stmtQuery,
				ExecContext:  &opts.stmtExecContext,
				QueryContext: &opts.stmtQueryContext,
				Rows: &rowsOptions{
					Close:         &opts.rowsClose,
					Next:          &opts.rowsNext,
					NextResultSet: &opts.rowsNextResultSet,
				},
			},
			ResetSession: &opts.connResetSession,
			Ping:         &opts.connPing,
			ExecContext:  &opts.connExecContext,
			QueryContext: &opts.connQueryContext,
			Rows: &rowsOptions{
				Close:         &opts.rowsClose,
				Next:          &opts.rowsNext,
				NextResultSet: &opts.rowsNextResultSet,
			},
		}

		buf := bytes.NewBuffer(nil)
		logger := slog.New(NewTextHandler(buf, nil))
		dw := wrapDriver(&mockErrorDiverContext{},
			newLogger(logger, opts),
			&driverOptions{
				IDGen:         opts.idGen,
				connIDKey:     opts.connIDKey,
				Open:          &opts.driverOpen,
				OpenConnector: &opts.driverOpenConnector,
				Conn:          connOptions,
				Connector: &connectorOptions{
					Connect: &opts.connectorConnect,
					Conn:    connOptions,
				},
			},
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
}
