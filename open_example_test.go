package sqlslog_test

import (
	"context"
	"log/slog"
	"os"

	sqlslog "github.com/akm/sql-slog"
	"github.com/akm/sql-slog/opts"
	// _ "github.com/mattn/go-sqlite3"
)

func ExampleOpen() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	db, err := sqlslog.Open(ctx, "sqlite3", dsn)
	if err != nil {
		// Handle error
	}
	defer db.Close()
	// Use db as a regular *sql.DB
}

func ExampleOpen_withLevel() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	logger := slog.New(opts.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: opts.LevelTrace,
	}))
	db, _ := sqlslog.Open(ctx, "sqlite3", dsn, opts.Logger(logger))
	defer db.Close()
}

func ExampleOpen_withStmtQueryContext() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	logger := slog.New(opts.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: opts.LevelTrace,
	}))
	db, _ := sqlslog.Open(ctx, "sqlite3", dsn,
		opts.Logger(logger),
		opts.StmtQueryContext(func(o *opts.StepOptions) {
			o.SetLevel(opts.LevelDebug)
		}),
	)
	defer db.Close()
}
