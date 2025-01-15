package sqlslog_test

import (
	"context"
	"log/slog"
	"os"

	sqlslog "github.com/akm/sql-slog"
	// _ "github.com/mattn/go-sqlite3"
)

func ExampleOpen() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db, err := sqlslog.Open(ctx, "sqlite3", dsn, logger)
	if err != nil {
		// Handle error
	}
	defer db.Close()

	// Use db as a regular *sql.DB
}
