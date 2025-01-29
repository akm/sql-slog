package opts_test

import (
	"context"
	"log/slog"
	"os"

	sqlslog "github.com/akm/sql-slog"
	"github.com/akm/sql-slog/opts"
)

func ExampleNewJSONHandler() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	logger := slog.New(opts.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	db, _ := sqlslog.Open(ctx, "sqlite3", dsn, opts.Logger(logger))
	defer db.Close()
}

func ExampleNewTextHandler() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	logger := slog.New(opts.NewTextHandler(os.Stdout, nil))
	db, _ := sqlslog.Open(ctx, "sqlite3", dsn, opts.Logger(logger))
	defer db.Close()
}
