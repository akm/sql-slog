package sqlslog

import (
	"bytes"
	"context"
	"database/sql/driver"
	"errors"
	"log/slog"
	"testing"
)

type mockStmtForWrapStmt struct {
	error error
}

var _ driver.Stmt = (*mockStmtForWrapStmt)(nil)

// Close implements driver.Stmt.
func (m *mockStmtForWrapStmt) Close() error {
	return m.error
}

// Exec implements driver.Stmt.
func (m *mockStmtForWrapStmt) Exec([]driver.Value) (driver.Result, error) {
	return nil, m.error
}

// NumInput implements driver.Stmt.
func (m *mockStmtForWrapStmt) NumInput() int {
	panic("unimplemented")
}

// Query implements driver.Stmt.
func (m *mockStmtForWrapStmt) Query([]driver.Value) (driver.Rows, error) {
	return nil, m.error
}

func TestWrapStmt(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		if wrapStmt(nil, nil) != nil {
			t.Fatal("Expected nil")
		}
	})
	t.Run("implements driver.Stmt but not stmtWithContext", func(t *testing.T) {
		t.Parallel()
		mock := &mockStmtForWrapStmt{}
		logger := &stepLogger{}
		stmt := wrapStmt(mock, logger)
		if stmt == nil {
			t.Fatal("Expected non-nil")
		}
	})

	t.Run("Query", func(t *testing.T) {
		t.Parallel()
		dummyError := errors.New("unexpected Query error")
		mock := &mockStmtForWrapStmt{
			error: dummyError,
		}

		buf := bytes.NewBuffer(nil)
		logger := slog.New(NewJSONHandler(buf, nil))
		wrapped := wrapStmt(mock, newStepLogger(logger, newOptions("dummy")))
		_, err := wrapped.Query(nil) // nolint:staticcheck
		if err == nil {
			t.Fatal("Expected non-nil")
		}
		if !errors.Is(err, dummyError) {
			t.Fatalf("Expected %q but got %q", dummyError, err)
		}
	})
}

type mockErrorStmtWithContext struct {
	mockStmtForWrapStmt
	error error
}

var (
	_ driver.Stmt             = (*mockErrorStmtWithContext)(nil)
	_ driver.StmtExecContext  = (*mockErrorStmtWithContext)(nil)
	_ driver.StmtQueryContext = (*mockErrorStmtWithContext)(nil)
)

func (m *mockErrorStmtWithContext) QueryContext(context.Context, []driver.NamedValue) (driver.Rows, error) {
	return nil, m.error
}

func (m *mockErrorStmtWithContext) ExecContext(context.Context, []driver.NamedValue) (driver.Result, error) {
	return nil, m.error
}

func TestWithMockErrorStmtWithContext(t *testing.T) {
	t.Parallel()
	dummyError := errors.New("unexpected QueryContext error")
	mock := &mockErrorStmtWithContext{
		mockStmtForWrapStmt: mockStmtForWrapStmt{},
		error:               dummyError,
	}

	buf := bytes.NewBuffer(nil)
	logger := slog.New(NewJSONHandler(buf, nil))
	wrapped := wrapStmt(mock, newStepLogger(logger, newOptions("dummy")))
	stmtWithQueryContext, ok := wrapped.(driver.StmtQueryContext)
	if !ok {
		t.Fatal("Expected StmtQueryContext")
	}
	_, err := stmtWithQueryContext.QueryContext(context.TODO(), nil)
	if err == nil {
		t.Fatal("Expected non-nil")
	}
	if !errors.Is(err, dummyError) {
		t.Fatalf("Expected %q but got %q", dummyError, err)
	}
}
