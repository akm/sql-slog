package sqlslog

import (
	"database/sql/driver"
	"fmt"
	"log/slog"
	"testing"
)

func TestWrapRows(t *testing.T) {
	if wrapRows(nil, nil) != nil {
		t.Fatal("Expected nil")
	}
}

type mockRows struct{}

var _ driver.Rows = (*mockRows)(nil)

// Close implements driver.Rows.
func (m *mockRows) Close() error {
	panic("unimplemented")
}

// Columns implements driver.Rows.
func (m *mockRows) Columns() []string {
	panic("unimplemented")
}

// Next implements driver.Rows.
func (m *mockRows) Next(dest []driver.Value) error {
	panic("unimplemented")
}

func TestWithMockRows(t *testing.T) {
	wrapped := &rowsWrapper{original: &mockRows{}, logger: newLogger(slog.Default(), nil)}
	t.Run("ColumnTypeScanType", func(t *testing.T) {
		res := wrapped.ColumnTypeScanType(0)
		if res == nil {
			t.Fatal("Expected non-nil")
		}
	})
	t.Run("ColumnTypeDatabaseTypeName", func(t *testing.T) {
		res := wrapped.ColumnTypeDatabaseTypeName(0)
		if res != "" {
			t.Fatal("Expected empty")
		}
	})
}

func TestHandleRowsNextError(t *testing.T) {
	complete, attrs := HandleRowsNextError(fmt.Errorf("dummy"))
	if complete {
		t.Fatal("Expected false")
	}
	if attrs != nil {
		t.Fatal("Expected nil")
	}
}
