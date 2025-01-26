package sqlslog

import (
	"database/sql/driver"
	"testing"
)

type mockStmtForWrapStmt struct {
}

var _ driver.Stmt = (*mockStmtForWrapStmt)(nil)

// Close implements driver.Stmt.
func (m *mockStmtForWrapStmt) Close() error {
	panic("unimplemented")
}

// Exec implements driver.Stmt.
func (m *mockStmtForWrapStmt) Exec([]driver.Value) (driver.Result, error) {
	panic("unimplemented")
}

// NumInput implements driver.Stmt.
func (m *mockStmtForWrapStmt) NumInput() int {
	panic("unimplemented")
}

// Query implements driver.Stmt.
func (m *mockStmtForWrapStmt) Query([]driver.Value) (driver.Rows, error) {
	panic("unimplemented")
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
		logger := &logger{}
		stmt := wrapStmt(mock, logger)
		if stmt == nil {
			t.Fatal("Expected non-nil")
		}
	})
}
