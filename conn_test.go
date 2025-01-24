package sqlslog

import (
	"context"
	"database/sql/driver"
	"fmt"
	"log/slog"
	"testing"
)

type mockConnForWrapConn struct {
}

// Begin implements driver.Conn.
func (m *mockConnForWrapConn) Begin() (driver.Tx, error) {
	panic("unimplemented")
}

// Close implements driver.Conn.
func (m *mockConnForWrapConn) Close() error {
	panic("unimplemented")
}

// Prepare implements driver.Conn.
func (m *mockConnForWrapConn) Prepare(query string) (driver.Stmt, error) {
	panic("unimplemented")
}

var _ driver.Conn = (*mockConnForWrapConn)(nil)

func TestWrapConn(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		if wrapConn(nil, nil) != nil {
			t.Fatal("Expected nil")
		}
	})
	t.Run("implements driver.Conn but not connWithContext", func(t *testing.T) {
		mock := &mockConnForWrapConn{}
		logger := &logger{}
		conn := wrapConn(mock, logger)
		if conn == nil {
			t.Fatal("Expected non-nil")
		}
	})
}

func TestConnExecContextErrorHandler(t *testing.T) {
	errHandler := ConnExecContextErrorHandler("mysql")
	complete, attrs := errHandler(fmt.Errorf("dummy"))
	if complete {
		t.Fatal("Expected false")
	}
	if attrs != nil {
		t.Fatal("Expected nil")
	}
}

func TestConnQueryContextErrorHandler(t *testing.T) {
	t.Run("mysql", func(t *testing.T) {
		errHandler := ConnQueryContextErrorHandler("mysql")
		t.Run("nil error", func(t *testing.T) {
			complete, attrs := errHandler(nil)
			if !complete {
				t.Fatal("Expected true")
			}
			if attrs != nil {
				t.Fatal("Expected nil")
			}
		})
		t.Run("unexpected error", func(t *testing.T) {
			complete, attrs := errHandler(fmt.Errorf("dummy"))
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
func (m *mockErrorConn) Prepare(query string) (driver.Stmt, error) {
	return nil, m.error
}

// BeginTx implements driver.ConnBeginTx.
func (m *mockErrorConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return nil, m.error
}

// PrepareContext implements driver.ConnPrepareContext.
func (m *mockErrorConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return nil, m.error
}

// QueryContext implements driver.QueryerContext.
func (m *mockErrorConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return nil, m.error
}

// ExecContext implements driver.ExecerContext.
func (m *mockErrorConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return nil, m.error
}

var _ driver.Conn = (*mockErrorConn)(nil)
var _ driver.ConnBeginTx = (*mockErrorConn)(nil)
var _ driver.ConnPrepareContext = (*mockErrorConn)(nil)
var _ driver.ExecerContext = (*mockErrorConn)(nil)
var _ driver.QueryerContext = (*mockErrorConn)(nil)

// var _ driver.Pinger = (*mockErrorConn)(nil) // not implemented for the test below

func TestWithMockErrorConn(t *testing.T) {
	logger := newLogger(slog.Default(), newOptions("sqlite3"))
	w := wrapConn(newMockErrConn(fmt.Errorf("unexpected error")), logger)
	t.Run("Begin", func(t *testing.T) {
		if _, err := w.Begin(); err == nil { //nolint:staticcheck
			t.Fatal("Expected error")
		}
	})
	t.Run("Prepare", func(t *testing.T) {
		if _, err := w.Prepare("dummy"); err == nil {
			t.Fatal("Expected error")
		}
	})
	t.Run("BeginTx", func(t *testing.T) {
		if _, err := w.(driver.ConnBeginTx).BeginTx(context.Background(), driver.TxOptions{}); err == nil {
			t.Fatal("Expected error")
		}
	})
	t.Run("PrepareContext", func(t *testing.T) {
		if _, err := w.(driver.ConnPrepareContext).PrepareContext(context.Background(), "dummy"); err == nil {
			t.Fatal("Expected error")
		}
	})
}

func TestPingInCase(t *testing.T) {
	logger := newLogger(slog.Default(), newOptions("sqlite3"))
	conn := newMockErrConn(nil)
	w := &connWithContextWrapper{
		connWrapper: connWrapper{
			original: conn,
			logger:   logger,
		},
		originalConn: conn,
	}
	if err := w.Ping(context.Background()); err != nil {
		t.Fatal("Unexpected error")
	}
}
