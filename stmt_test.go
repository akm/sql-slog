package sqlslog

import (
	"testing"
)

func TestWrapStmt(t *testing.T) {
	if wrapStmt(nil, nil) != nil {
		t.Fatal("Expected nil")
	}
}
