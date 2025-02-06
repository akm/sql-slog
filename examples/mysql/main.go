package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"slices"
	"time"

	sqlslog "github.com/akm/sql-slog"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var logLevel sqlslog.Level
	if slices.Contains(os.Args, "debug") {
		logLevel = sqlslog.LevelDebug
	} else if slices.Contains(os.Args, "trace") {
		logLevel = sqlslog.LevelTrace
	} else if slices.Contains(os.Args, "verbose") {
		logLevel = sqlslog.LevelVerbose
	} else {
		logLevel = sqlslog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: logLevel}

	var handler slog.Handler
	if slices.Contains(os.Args, "json") {
		handler = sqlslog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = sqlslog.NewTextHandler(os.Stdout, opts)
	}
	logger := slog.New(handler)

	dbName := "app1"
	dbPort := 3306
	dsn := fmt.Sprintf("root@tcp(localhost:%d)/%s", dbPort, dbName)

	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DATABASE", dbName)
	if err := exec.Command("docker", "compose", "-f", "docker-compose.yml", "up", "-d").Run(); err != nil {
		panic(err)
	}
	defer func() {
		if err := exec.Command("docker", "compose", "-f", "docker-compose.yml", "down").Run(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	// Open a database
	db, err := sqlslog.Open(ctx, "mysql", dsn,
		sqlslog.Logger(logger),
		sqlslog.ConnQueryContext(func(o *sqlslog.StepOptions) {
			o.SetLevel(sqlslog.LevelDebug)
		}),
	)
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	for i := 0; i < 10; i++ {
		if err := db.PingContext(ctx); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	// Create a table
	query := "CREATE TABLE IF NOT EXISTS test1 (id INT PRIMARY KEY, name VARCHAR(255))"
	if _, err := db.ExecContext(ctx, query); err != nil {
		logger.Error("Failed to create a table", "error", err)
		os.Exit(1)
	}

	for i := 0; i < 10; i++ {
		if err := db.PingContext(ctx); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	// Insert a record in a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Failed to begin a transaction", "error", err)
		os.Exit(1)
	}

	query = "INSERT INTO test1 (id, name) VALUES (?, ?)"
	if _, err := tx.ExecContext(ctx, query, 1, "Alice"); err != nil {
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
