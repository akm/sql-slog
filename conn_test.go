package sqlslog

import (
	"database/sql/driver"
	"fmt"
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
