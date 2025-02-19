package sqlslog_test

import (
	"context"

	sqlslog "github.com/akm/sql-slog"
	// _ "github.com/mattn/go-sqlite3"
)

func ExampleOpen() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	db, logger, err := sqlslog.Open(ctx, "sqlite3", dsn)
	if err != nil {
		// Handle error
	}
	defer db.Close()
	// Use db as a regular *sql.DB
	logger.InfoContext(ctx, "Hello, World!")
}

func ExampleOpen_withLevel() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	db, _, _ := sqlslog.Open(ctx, "sqlite3", dsn,
		sqlslog.LogLevel(sqlslog.LevelTrace),
	)
	defer db.Close()
}

func ExampleOpen_withStmtQueryContext() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	db, _, _ := sqlslog.Open(ctx, "sqlite3", dsn,
		sqlslog.LogLevel(sqlslog.LevelTrace),
		sqlslog.StmtQueryContext(func(o *sqlslog.StepOptions) {
			o.SetLevel(sqlslog.LevelDebug)
		}),
	)
	defer db.Close()
}
