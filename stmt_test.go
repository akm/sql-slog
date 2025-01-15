package sqlslog

import (
	"log/slog"
	"testing"
)

func TestWrapStmt(t *testing.T) {
	if wrapStmt(nil, slog.Default()) != nil {
		t.Fatal("Expected nil")
	}
}
