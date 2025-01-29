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

func ExampleLogger() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db, _ := sqlslog.Open(ctx, "sqlite3", dsn, opts.Logger(logger))
	defer db.Close()
}

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

func ExampleSetStepLogMsgFormatter() {
	opts.SetStepLogMsgFormatter(func(name string, event opts.Event) string {
		return name + "/" + event.String()
	})
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	logger := slog.New(opts.NewJSONHandler(os.Stdout, nil))
	db, _ := sqlslog.Open(ctx, "sqlite3", dsn, opts.Logger(logger))
	defer db.Close()
}
