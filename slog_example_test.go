package sqlslog_test

import (
	"context"
	"log/slog"
	"os"

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

func ExampleHandler() {
	// Normal slog.Handler with sqlslog.ReplaceLevelAttr knows how to show LevelTrace and LevelVerbose.
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: sqlslog.MergeReplaceAttrs(
			sqlslog.ReplaceLevelAttr,
			removeTimeAndDuration, // for testing
		),
		Level: sqlslog.LevelVerbose,
	})

	dsn := "dummy-dsn"
	ctx := context.TODO()

	logger := sqlslog.New("mock", dsn, sqlslog.Handler(handler)).Logger()
	logger.Log(ctx, slog.Level(sqlslog.LevelTrace), "Foo")
	logger.Log(ctx, slog.Level(sqlslog.LevelVerbose), "Bar")

	// Output:
	// level=TRACE msg=Foo
	// level=VERBOSE msg=Bar
}

func ExampleHandler_withoutReplaceLevelAttr() {
	// Normal slog.Handler without sqlslog.ReplaceLevelAttr does not know how to show LevelTrace and LevelVerbose.
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: removeTimeAndDuration,
		Level:       sqlslog.LevelVerbose,
	})

	dsn := "dummy-dsn"
	ctx := context.TODO()
	logger := sqlslog.New("mock", dsn, sqlslog.Handler(handler)).Logger()

	logger.Log(ctx, slog.Level(sqlslog.LevelTrace), "Foo")
	logger.Log(ctx, slog.Level(sqlslog.LevelVerbose), "Bar")

	// Output:
	// level=DEBUG-4 msg=Foo
	// level=DEBUG-8 msg=Bar
}
