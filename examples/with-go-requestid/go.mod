module example-with-go-requestid

go 1.23.2

require (
	github.com/akm/go-requestid v0.3.1
	github.com/akm/sql-slog v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.14.24
)

require github.com/akm/slogctx v0.5.1 // indirect

replace github.com/akm/sql-slog => ../..
