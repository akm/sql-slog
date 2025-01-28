package sqlslog

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"log/slog"
	"testing"
)

func TestWrapRows(t *testing.T) {
	t.Parallel()
	if WrapRows(nil, nil) != nil {
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
func (m *mockRows) Next([]driver.Value) error {
	panic("unimplemented")
}

func TestWithMockRows(t *testing.T) {
	t.Parallel()
	wrapped := &rowsWrapper{original: &mockRows{}, logger: newLogger(slog.Default(), nil)}
	t.Run("ColumnTypeScanType", func(t *testing.T) {
		t.Parallel()
		res := wrapped.ColumnTypeScanType(0)
		if res == nil {
			t.Fatal("Expected non-nil")
		}
	})
	t.Run("ColumnTypeDatabaseTypeName", func(t *testing.T) {
		t.Parallel()
		res := wrapped.ColumnTypeDatabaseTypeName(0)
		if res != "" {
			t.Fatal("Expected empty")
		}
	})
}

type mockRowsNextResultSet struct {
	mockRows
	error error
}

func (m *mockRowsNextResultSet) Close() error {
	return m.error
}

func (m *mockRowsNextResultSet) Columns() []string {
	panic("unimplemented")
}

func (m *mockRowsNextResultSet) HasNextResultSet() bool {
	panic("unimplemented")
}

func (m *mockRowsNextResultSet) Next([]driver.Value) error {
	return m.error
}

func (m *mockRowsNextResultSet) NextResultSet() error {
	return m.error
}

var _ driver.RowsNextResultSet = (*mockRowsNextResultSet)(nil)

func TestRowsNextResultSet(t *testing.T) {
	t.Parallel()
	errMsg := "unpected RNRS error"
	rows := &mockRowsNextResultSet{
		mockRows: mockRows{},
		error:    errors.New(errMsg),
	}
	buf := bytes.NewBuffer(nil)
	logger := slog.New(NewJSONHandler(buf, nil))
	wrapped := WrapRows(rows, newLogger(logger, newOptions("dummy")))
	wrappedRNRS, ok := wrapped.(driver.RowsNextResultSet)
	if !ok {
		t.Fatal("Expected true")
	}
	err := wrappedRNRS.NextResultSet()
	if err == nil {
		t.Fatal("Expected non-nil")
	}
	if err.Error() != errMsg {
		t.Fatalf("Expected %q, got %q", errMsg, err.Error())
	}
}

func TestHandleRowsNextError(t *testing.T) {
	t.Parallel()
	complete, attrs := HandleRowsNextError(errors.New("dummy"))
	if complete {
		t.Fatal("Expected false")
	}
	if attrs != nil {
		t.Fatal("Expected nil")
	}
}
