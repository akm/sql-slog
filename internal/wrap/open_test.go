package wrap

import (
	"database/sql/driver"
	"errors"
	"log/slog"
	"testing"
)

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

func buildOpenOptions(options *Options) *OpenOptions {
	connOptions := &ConnOptions{
		IDGen:   options.IDGen,
		Begin:   &options.ConnBegin,
		BeginTx: &options.ConnBeginTx,
		TxIDKey: options.TxIDKey,
		Tx: &TxOptions{
			Commit:   &options.TxCommit,
			Rollback: &options.TxRollback,
		},
		Close:          &options.ConnClose,
		Prepare:        &options.ConnPrepare,
		PrepareContext: &options.ConnPrepareContext,
		StmtIDKey:      options.StmtIDKey,
		Stmt: &StmtOptions{
			Close:        &options.StmtClose,
			Exec:         &options.StmtExec,
			Query:        &options.StmtQuery,
			ExecContext:  &options.StmtExecContext,
			QueryContext: &options.StmtQueryContext,
			Rows: &RowsOptions{
				Close:         &options.RowsClose,
				Next:          &options.RowsNext,
				NextResultSet: &options.RowsNextResultSet,
			},
		},
		ResetSession: &options.ConnResetSession,
		Ping:         &options.ConnPing,
		ExecContext:  &options.ConnExecContext,
		QueryContext: &options.ConnQueryContext,
		Rows: &RowsOptions{
			Close:         &options.RowsClose,
			Next:          &options.RowsNext,
			NextResultSet: &options.RowsNextResultSet,
		},
	}
	return &OpenOptions{
		Open: &options.SqlslogOpen,
		Driver: &DriverOptions{
			IDGen:         options.IDGen,
			ConnIDKey:     options.ConnIDKey,
			Open:          &options.DriverOpen,
			OpenConnector: &options.DriverOpenConnector,
			Conn:          connOptions,
			Connector: &ConnectorOptions{
				Connect: &options.ConnectorConnect,
				Conn:    connOptions,
			},
		},
	}
}
