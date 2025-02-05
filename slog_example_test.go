package sqlslog_test

import (
	"context"
	"log/slog"
	"os"

	sqlslog "github.com/akm/sql-slog"
	// _ "github.com/mattn/go-sqlite3"
)

func ExampleNewJSONHandler() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	logger := slog.New(sqlslog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	db, _ := sqlslog.Open(ctx, "sqlite3", dsn, sqlslog.Logger(logger))
	defer db.Close()
}

func ExampleNewTextHandler() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	logger := slog.New(sqlslog.NewTextHandler(os.Stdout, nil))
	db, _ := sqlslog.Open(ctx, "sqlite3", dsn, sqlslog.Logger(logger))
	defer db.Close()
}
