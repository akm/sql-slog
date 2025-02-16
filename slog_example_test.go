package sqlslog_test

import (
	"context"
	"log/slog"

	sqlslog "github.com/akm/sql-slog"
	// _ "github.com/mattn/go-sqlite3"
)

func ExampleNewJSONHandler() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	db, logger, _ := sqlslog.Open(ctx, "sqlite3", dsn,
		sqlslog.HandlerFunc(sqlslog.NewJSONHandler),
		sqlslog.HandlerOptions(&slog.HandlerOptions{Level: sqlslog.LevelDebug}),
	)
	defer db.Close()
	logger.InfoContext(ctx, "Hello, World!")
}

func ExampleNewTextHandler() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	db, logger, _ := sqlslog.Open(ctx, "sqlite3", dsn,
		sqlslog.HandlerFunc(sqlslog.NewTextHandler),
	)
	defer db.Close()
	logger.InfoContext(ctx, "Hello, World!")
}
