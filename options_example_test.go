package sqlslog_test

import (
	"context"
	"log/slog"
	"os"

	sqlslog "github.com/akm/sql-slog"
	// _ "github.com/mattn/go-sqlite3"
)

func ExampleLogger() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db, _ := sqlslog.Open(ctx, "sqlite3", dsn, sqlslog.Logger(logger))
	defer db.Close()
}

func ExampleSetStepLogMsgFormatter() {
	sqlslog.SetStepLogMsgFormatter(func(name string, event sqlslog.Event) string {
		return name + "/" + event.String()
	})
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	logger := slog.New(sqlslog.NewJSONHandler(os.Stdout, nil))
	db, _ := sqlslog.Open(ctx, "sqlite3", dsn, sqlslog.Logger(logger))
	defer db.Close()
}
