package sqlslog

import (
	"log/slog"
	"testing"
)

func TestWrapConn(t *testing.T) {
	if wrapConn(nil, slog.Default()) != nil {
		t.Fatal("Expected nil")
	}
}
