package sqlslog_test

import (
	"context"

	sqlslog "github.com/akm/sql-slog"
	// _ "github.com/mattn/go-sqlite3"
)

func ExampleHandlerFunc() {
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	db, logger, _ := sqlslog.Open(ctx, "sqlite3", dsn,
		sqlslog.HandlerFunc(sqlslog.NewJSONHandler),
	)
	defer db.Close()
	logger.InfoContext(ctx, "Hello, World!")
}

func ExampleSetStepEventMsgBuilder() {
	sqlslog.SetStepEventMsgBuilder(func(step sqlslog.Step, event sqlslog.Event) string {
		return step.String() + "/" + event.String()
	})
	dsn := "file::memory:?cache=shared"
	ctx := context.TODO()
	db, _, _ := sqlslog.Open(ctx, "sqlite3", dsn,
		sqlslog.HandlerFunc(sqlslog.NewJSONHandler),
	)
	defer db.Close()
}
