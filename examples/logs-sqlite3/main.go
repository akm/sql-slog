package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"slices"

	sqlslog "github.com/akm/sql-slog"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logLevel := sqlslog.ParseLevelWithDefault(os.Args[1], sqlslog.LevelInfo)

	var handlerFunc func(io.Writer, *slog.HandlerOptions) slog.Handler
	if slices.Contains(os.Args, "json") {
		handlerFunc = sqlslog.NewJSONHandler
	} else {
		handlerFunc = sqlslog.NewTextHandler
	}

	ctx := context.Background()
	dsn := "file::memory:?cache=shared"

	// Open a database
	db, logger, err := sqlslog.Open(ctx, "sqlite3", dsn,
		sqlslog.LogLevel(logLevel),
		sqlslog.HandlerFunc(handlerFunc),
		sqlslog.ConnQueryContext(func(o *sqlslog.StepOptions) {
			o.SetLevel(sqlslog.LevelDebug)
		}),
	)
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create a table
	query := "CREATE TABLE IF NOT EXISTS test1 (id INTEGER PRIMARY KEY, name VARCHAR(255))"
	if _, err := db.ExecContext(ctx, query); err != nil {
		logger.Error("Failed to create a table", "error", err)
		os.Exit(1)
	}

	// Insert a record in a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Failed to begin a transaction", "error", err)
		os.Exit(1)
	}

	query = "INSERT INTO test1 (name) VALUES (?)"
	if _, err := tx.ExecContext(ctx, query, "Alice"); err != nil {
		logger.Error("Failed to insert a record", "error", err)
		if err := tx.Rollback(); err != nil {
			logger.Error("Failed to rollback a transaction", "error", err)
		}
		os.Exit(1)
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit a transaction", "error", err)
		os.Exit(1)
	}

	// Select records
	rows, err := db.QueryContext(ctx, "SELECT * FROM test1")
	if err != nil {
		logger.Error("Failed to select records", "error", err)
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			logger.Error("Failed to scan a record", "error", err)
			os.Exit(1)
		}
		logger.InfoContext(ctx, "Record", "id", id, "name", name)
	}
}
