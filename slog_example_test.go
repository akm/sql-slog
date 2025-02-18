package sqlslog_test

import (
	"context"
	"log/slog"

	sqlslog "github.com/akm/sql-slog"
	// _ "github.com/mattn/go-sqlite3"
)

func removeTimeAndDuration(groups []string, a slog.Attr) slog.Attr {
	if len(groups) == 0 {
		switch a.Key {
		case slog.TimeKey:
			return slog.Attr{}
		case "duration":
			return slog.Attr{}
		}
	}
	return a
}

func ExampleNewJSONHandler() {
	dsn := "dummy-dsn"
	ctx := context.TODO()
	db, logger, _ := sqlslog.Open(ctx, "mock", dsn,
		sqlslog.HandlerFunc(sqlslog.NewJSONHandler),
		sqlslog.LogReplaceAttr(removeTimeAndDuration),
	)
	defer db.Close()
	logger.InfoContext(ctx, "Hello, World!")

	// Output:
	// {"level":"INFO","msg":"Open","driver":"mock","dsn":"dummy-dsn"}
	// {"level":"INFO","msg":"Hello, World!"}
}

func ExampleNewTextHandler() {
	dsn := "dummy-dsn"
	ctx := context.TODO()
	db, logger, _ := sqlslog.Open(ctx, "mock", dsn,
		sqlslog.HandlerFunc(sqlslog.NewTextHandler),
		sqlslog.LogReplaceAttr(removeTimeAndDuration),
	)
	defer db.Close()
	logger.InfoContext(ctx, "Hello, World!")

	// Output:
	// level=INFO msg=Open driver=mock dsn=dummy-dsn
	// level=INFO msg="Hello, World!"
}
