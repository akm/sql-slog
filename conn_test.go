package sqlslog

import (
	"context"
	"database/sql/driver"
	"errors"
	"log/slog"
	"testing"
)

type mockConnForWrapConn struct{}

// Begin implements driver.Conn.
func (m *mockConnForWrapConn) Begin() (driver.Tx, error) {
	panic("unimplemented")
}

// Close implements driver.Conn.
func (m *mockConnForWrapConn) Close() error {
	panic("unimplemented")
}

// Prepare implements driver.Conn.
func (m *mockConnForWrapConn) Prepare(string) (driver.Stmt, error) {
	panic("unimplemented")
}

var _ driver.Conn = (*mockConnForWrapConn)(nil)

func TestWrapConn(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		if wrapConn(nil, nil, nil) != nil {
			t.Fatal("Expected nil")
		}
	})
	t.Run("implements driver.Conn but not connWithContext", func(t *testing.T) {
		t.Parallel()
		mock := &mockConnForWrapConn{}
		logger := &stepLogger{}
		connOptions := defaultConnOptions("dummy", StepEventMsgWithoutEventName)
		conn := wrapConn(mock, logger, connOptions)
		if conn == nil {
			t.Fatal("Expected non-nil")
		}

		t.Run("skip wrapped driver.Conn object", func(t *testing.T) {
			res := wrapConn(conn, logger, connOptions)
			if res != conn {
				t.Fatal("Expected same object")
			}
		})
	})
}

func TestConnExecContextErrorHandler(t *testing.T) {
	t.Parallel()
	errHandler := ConnExecContextErrorHandler("mysql")
	complete, attrs := errHandler(errors.New("dummy"))
	if complete {
		t.Fatal("Expected false")
	}
	if attrs != nil {
		t.Fatal("Expected nil")
	}
}

func TestConnQueryContextErrorHandler(t *testing.T) {
	t.Parallel()
	t.Run("mysql", func(t *testing.T) {
		t.Parallel()
		errHandler := ConnQueryContextErrorHandler("mysql")
		t.Run("nil error", func(t *testing.T) {
			t.Parallel()
			complete, attrs := errHandler(nil)
			if !complete {
				t.Fatal("Expected true")
			}
			if attrs != nil {
				t.Fatal("Expected nil")
			}
		})
		t.Run("unexpected error", func(t *testing.T) {
			t.Parallel()
			complete, attrs := errHandler(errors.New("dummy"))
			if complete {
				t.Fatal("Expected false")
			}
			if attrs != nil {
				t.Fatal("Expected nil")
			}
		})
	})
}

type mockErrorConn struct {
	error error
}

func newMockErrConn(err error) *mockErrorConn {
	return &mockErrorConn{error: err}
}

// Begin implements driver.Conn.
func (m *mockErrorConn) Begin() (driver.Tx, error) {
	return nil, m.error
}

// Close implements driver.Conn.
func (m *mockErrorConn) Close() error {
	return m.error
}

// Prepare implements driver.Conn.
func (m *mockErrorConn) Prepare(string) (driver.Stmt, error) {
	return nil, m.error
}

// BeginTx implements driver.ConnBeginTx.
func (m *mockErrorConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return nil, m.error
}

// PrepareContext implements driver.ConnPrepareContext.
func (m *mockErrorConn) PrepareContext(context.Context, string) (driver.Stmt, error) {
	return nil, m.error
}

// QueryContext implements driver.QueryerContext.
func (m *mockErrorConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return nil, m.error
}

// ExecContext implements driver.ExecerContext.
func (m *mockErrorConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	return nil, m.error
}

var (
	_ driver.Conn               = (*mockErrorConn)(nil)
	_ driver.ConnBeginTx        = (*mockErrorConn)(nil)
	_ driver.ConnPrepareContext = (*mockErrorConn)(nil)
	_ driver.ExecerContext      = (*mockErrorConn)(nil)
	_ driver.QueryerContext     = (*mockErrorConn)(nil)
)

// var _ driver.Pinger = (*mockErrorConn)(nil) // not implemented for the test below

func TestWithMockErrorConn(t *testing.T) {
	t.Parallel()
	logger := newStepLogger(slog.Default(), defaultStepLoggerOptions())
	connOptions := defaultConnOptions("sqlite3", StepEventMsgWithoutEventName)
	w := wrapConn(newMockErrConn(errors.New("unexpected error")), logger, connOptions)
	t.Run("Begin", func(t *testing.T) {
		t.Parallel()
		if _, err := w.Begin(); err == nil { //nolint:staticcheck
			t.Fatal("Expected error")
		}
	})
	t.Run("Prepare", func(t *testing.T) {
		t.Parallel()
		if _, err := w.Prepare("dummy"); err == nil {
			t.Fatal("Expected error")
		}
	})
	t.Run("BeginTx", func(t *testing.T) {
		t.Parallel()
		if _, err := w.(driver.ConnBeginTx).BeginTx(context.Background(), driver.TxOptions{}); err == nil {
			t.Fatal("Expected error")
		}
	})
	t.Run("PrepareContext", func(t *testing.T) {
		t.Parallel()
		if _, err := w.(driver.ConnPrepareContext).PrepareContext(context.Background(), "dummy"); err == nil {
			t.Fatal("Expected error")
		}
	})
}

func TestPingInCase(t *testing.T) {
	t.Parallel()
	logger := newStepLogger(slog.Default(), defaultStepLoggerOptions())
	conn := newMockErrConn(nil)
	w := &connWithContextWrapper{
		connWrapper: connWrapper{
			original: conn,
			logger:   logger,
			options:  defaultConnOptions("sqlite3", StepEventMsgWithoutEventName),
		},
		originalConn: conn,
	}
	if err := w.Ping(context.Background()); err != nil {
		t.Fatal("Unexpected error")
	}
}
