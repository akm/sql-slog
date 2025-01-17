package sqlslog

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
)

func TestOpen(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	logger := slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	db, err := Open(ctx, "invalid-driver", "", Logger(logger))
	if err == nil {
		t.Fatal("Expected error")
	}
	if err.Error() != "sql: unknown driver \"invalid-driver\" (forgotten import?)" {
		t.Fatalf("Unexpected error: %v", err)
	}
	if db != nil {
		t.Fatal("Expected nil db")
	}
}
