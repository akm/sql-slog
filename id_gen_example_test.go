package sqlslog_test

import (
	"context"
	cryptorand "crypto/rand"
	"log/slog"
	mathrandv2 "math/rand/v2"
	"os"

	sqlslog "github.com/akm/sql-slog"
	// _ "github.com/mattn/go-sqlite3"
)

func ExampleIDGenerator() {
	idGen := sqlslog.RandIntIDGenerator(
		mathrandv2.Int,
		[]byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		8,
	)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db, _ := sqlslog.Open(context.TODO(), "sqlite3", "file::memory:?cache=shared",
		sqlslog.Logger(logger),
		sqlslog.IDGenerator(idGen),
	)
	defer db.Close()
}

func ExampleRandIntIDGenerator() {
	idGen := sqlslog.RandIntIDGenerator(
		mathrandv2.Int,
		[]byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		8,
	)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db, _ := sqlslog.Open(context.TODO(), "sqlite3", "file::memory:?cache=shared",
		sqlslog.Logger(logger),
		sqlslog.IDGenerator(idGen),
	)
	defer db.Close()
}

func ExampleIDGenErrorSuppressor() {
	idGen := sqlslog.IDGenErrorSuppressor(
		sqlslog.RandReadIDGenerator(
			cryptorand.Read,
			[]byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
			8,
		),
		func(error) string { return "recovered" },
	)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db, _ := sqlslog.Open(context.TODO(), "sqlite3", "file::memory:?cache=shared",
		sqlslog.Logger(logger),
		sqlslog.IDGenerator(idGen),
	)
	defer db.Close()
}
