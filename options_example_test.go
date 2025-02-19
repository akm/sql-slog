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
		return "PRFIX:" + step.String() + "/" + event.String() + ":SUFFIX"
	})
	defer sqlslog.SetStepEventMsgBuilder(sqlslog.StepEventMsgWithoutEventName)

	dsn := "dummy-dsn"
	ctx := context.TODO()
	db, logger, _ := sqlslog.Open(ctx, "mock", dsn,
		sqlslog.LogReplaceAttr(removeTimeAndDuration), // for testing
	)
	defer db.Close()
	logger.InfoContext(ctx, "Hello, World!")

	// Output:
	// level=INFO msg=PRFIX:Open/Complete:SUFFIX driver=mock dsn=dummy-dsn
	// level=INFO msg="Hello, World!"
}
