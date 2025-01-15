package sqlslog

import (
	"log/slog"
	"testing"
)

func TestWrapRows(t *testing.T) {
	if wrapRows(nil, slog.Default()) != nil {
		t.Fatal("Expected nil")
	}
}
