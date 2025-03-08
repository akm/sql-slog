# sql-slog

[![CI](https://github.com/akm/sql-slog/actions/workflows/ci.yml/badge.svg)](https://github.com/akm/sql-slog/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/akm/sql-slog/graph/badge.svg?token=9BcanbSLut)](https://codecov.io/github/akm/sql-slog)
[![Go Report Card](https://goreportcard.com/badge/github.com/akm/sql-slog)](https://goreportcard.com/report/github.com/akm/sql-slog)
[![Go project version](https://badge.fury.io/go/github.com%2Fakm%2Fsql-slog.svg)](https://badge.fury.io/go/github.com%2Fakm%2Fsql-slog)
[![Enabled Linters](https://img.shields.io/badge/dynamic/yaml?url=https%3A%2F%2Fraw.githubusercontent.com%2Fakm%2Fsql-slog%2Frefs%2Fheads%2Fmain%2F.project.yaml&query=%24.linters&label=enabled%20linters&color=%2317AFC2)](.golangci.yml)
[![Documentation](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/akm/sql-slog)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/akm/sql-slog)](./go.mod)
[![license](https://img.shields.io/github/license/akm/sql-slog)](./LICENSE)

A logger for Go SQL database drivers with [log/slog](https://pkg.go.dev/log/slog) without modifying existing `*sql.DB` stdlib usage.

## FEATURES

- [x] Keep using (or re-use existing) `*sql.DB` as is.
- [x] No logger adapters. Just use [log/slog](https://pkg.go.dev/log/slog)
- [x] No dependencies
- [x] Leveled, detailed and configurable logging.
- [x] Duration
- [x] Trackable log output
- [x] Coverage 100%

See [godoc](https://pkg.go.dev/github.com/akm/sql-slog) for more details.

## LOG EXAMPLES

- [sqlite3](./examples/logs-sqlite3/results)
- [postgres](./examples/logs-postgres/results)
- [mysql](./examples/logs-mysql/results)

## INSTALL

To install sql-slog, use the following command:

```sh
go get -u github.com/akm/sql-slog
```

## USAGE

To use sql-slog, you can open a database connection with logging enabled as follows:

```golang
db, logger, err := sqlslog.Open(ctx, "mysql", dsn)
```

This is the easiest way to use sqlslog. It's similar to the usage of `Open` from `database/sql` like this:

```golang
db, err := sql.Open("mysql", dsn)
```

The differences are:

1. Pass `context.Context` as the first argument.
2. `*slog.Logger` is returned as the second argument.
3. `sqlslog.Open` can take a lot of [Option](https://pkg.go.dev/github.com/akm/sql-slog#Option).

See [godoc examples](https://pkg.go.dev/github.com/akm/sql-slog#example-Open) for more details.

## EXAMPLES

### [examples/with-go-requestid](./examples/with-go-requestid/)

An example showing how sql-slog works with [go-requestid](https://github.com/akm/go-requestid).
You can see DB query logs with request IDs in the same log like the following:

> time=2025-02-27T23:53:48.982+09:00 level=DEBUG msg=Conn.QueryContext conn_id=L1snTUaknlmsin8b query="SELECT id, title, status FROM todos" args=[] req_id=0JKGwDLjw77BjBnf

`conn_id` is a tracking ID for DB connection by sql-slog, and `req_id` is a tracking ID for HTTP request by go-requestid.

See [server-logs.txt](./examples/with-go-requestid/server-logs.txt) and [main.go](./examples/with-go-requestid/main.go) for more details.

## TEST

- [For MySQL](https://github.com/akm/sql-slog/blob/3f72cc68aefa9ac05b031d865dbdaec8a361c2c9/tests/mysql/low_level_with_context_test.go) for more details.
- [For PostgreSQL](https://github.com/akm/sql-slog/blob/3f72cc68aefa9ac05b031d865dbdaec8a361c2c9/tests/postgres/low_level_with_context_test.go) for more details.
- [For SQLite3](https://github.com/akm/sql-slog/blob/3f72cc68aefa9ac05b031d865dbdaec8a361c2c9/tests/sqlite3/low_level_without_context_test.go) for more details.

## MOTIVATION

I want to:

- Keep using `*sql.DB`.
  - To work with thin ORM
- Use log/slog
  - Leverage structured logging
  - Fetch and log `context.Context` values if needed

## REFERENCES

- [Stdlib sql.DB](https://github.com/golang/go/blob/master/src/database/sql/sql.go)
- [SQL driver interfaces](https://github.com/golang/go/blob/master/src/database/sql/driver/driver.go)
- [SQL driver implementation](https://go.dev/wiki/SQLDrivers)
- [log/slog](https://pkg.go.dev/log/slog)
- [Structured Logging with slog](https://go.dev/blog/slog)

## CONTRIBUTING

If you find a bug, typo, incorrect test, have an idea, or want to help with an existing issue, please create an issue or pull request.

## INSPIRED BY

- [github.com/simukti/sqldb-logger](https://github.com/simukti/sqldb-logger).

## LICENSE

[MIT](./LICENSE)
