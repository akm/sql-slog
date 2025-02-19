package sqlslog_test

import (
	"context"
	cryptorand "crypto/rand"
	mathrandv2 "math/rand/v2"

	sqlslog "github.com/akm/sql-slog"
	// _ "github.com/mattn/go-sqlite3"
)

func ExampleIDGenerator() {
	idGen := sqlslog.RandIntIDGenerator(
		mathrandv2.Int,
		[]byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		8,
	)

	db, _, _ := sqlslog.Open(context.TODO(), "sqlite3", "file::memory:?cache=shared",
		sqlslog.HandlerFunc(sqlslog.NewJSONHandler),
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

	db, _, _ := sqlslog.Open(context.TODO(), "sqlite3", "file::memory:?cache=shared",
		sqlslog.HandlerFunc(sqlslog.NewJSONHandler),
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

	db, _, _ := sqlslog.Open(context.TODO(), "sqlite3", "file::memory:?cache=shared",
		sqlslog.HandlerFunc(sqlslog.NewJSONHandler),
		sqlslog.IDGenerator(idGen),
	)
	defer db.Close()
}
