package opts_test

import (
	"context"
	"log/slog"
	"os"

	sqlslog "github.com/akm/sql-slog"
	"github.com/akm/sql-slog/opts"
)

func ExampleLogger() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
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
