# sql-slog

[![CI](https://github.com/akm/sql-slog/actions/workflows/ci.yml/badge.svg)](https://github.com/akm/sql-slog/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/akm/sql-slog/graph/badge.svg?token=9BcanbSLut)](https://codecov.io/github/akm/sql-slog)
[![Enabled Linters](https://img.shields.io/badge/dynamic/yaml?url=https%3A%2F%2Fraw.githubusercontent.com%2Fakm%2Fsql-slog%2Frefs%2Fheads%2Fmain%2F.project.yaml&query=%24.linters&label=enabled%20linters&color=%2317AFC2)](.golangci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/akm/sql-slog)](https://goreportcard.com/report/github.com/akm/sql-slog)
[![Documentation](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/akm/sql-slog)
![license](https://img.shields.io/github/license/akm/sql-slog)

A logger for Go SQL database drivers without modifying existing `*sql.DB` stdlib usage.

## LOG EXAMPLES

- [mysql](./examples/logs-mysql/results)
- [postgres](./examples/logs-postgres/results)
- [sqlite3](./examples/logs-sqlite3/results)

## FEATURES

- Keep using (or re-use existing) `*sql.DB` as is.
- No logger adapters. Just use [log/slog](https://pkg.go.dev/log/slog)
- No dependencies except stdlib.
- Leveled, detailed and configurable logging.
- Duration
- Trackable log output
  - conn_id
  - tx_id
  - stmt_id

## INSTALL

```sh
go get -u github.com/akm/sql-slog
```

## USAGE

This is a simple example of how to use `sql.Open`:

```golang
db, err := sql.Open("mysql", dsn)
```

When using sqlslog, you can use it like this:

```golang
ctx := context.TODO() // or a context.Context
db, err := sqlslog.Open(ctx, "mysql", dsn)
```

1. Replace `sql.Open` with `sqlslog.Open`.
2. Insert a `context.Context` as the first argument.

## Features

### Additional Log Levels

sqlslog provides additional log levels `LevelTrace` and `LevelVerbose` as [sqlslog.Level](https://pkg.go.dev/github.com/akm/sql-slog#Level).

To display the log levels correctly, the logger handler must be customized. You can create a handler using [sqlslog.NewJSONHandler](https://pkg.go.dev/github.com/akm/sql-slog#NewJSONHandler) and [sqlslog.NewTextHandler](https://pkg.go.dev/github.com/akm/sql-slog#NewTextHandler).

Pass an [sqlslog.Option](https://pkg.go.dev/github.com/akm/sql-slog#Option) created by [sqlslog.Logger](https://pkg.go.dev/github.com/akm/sql-slog#Logger) to [sqlslog.Open](https://pkg.go.dev/github.com/akm/sql-slog#Open) to use them.

```golang
db, err := sqlslog.Open(ctx, "sqlite3", dsn, sqlslog.Logger(yourLogger))
```

### Configurable Log Messages and Log Levels for Each Step

In sqlslog terms, a step is a logical operation in the database driver, such as a query, a ping, a prepare, etc.

A step has three events: start, error, and complete.

sqlslog provides a way to customize the log message and log level for each step event.

You can customize them using functions that take [StepOptions](https://pkg.go.dev/github.com/akm/sql-slog#StepOptions) and return [Option](https://pkg.go.dev/github.com/akm/sql-slog#Option), like [ConnPrepareContext](https://pkg.go.dev/github.com/akm/sql-slog#ConnPrepareContext) or [StmtQueryContext](https://pkg.go.dev/github.com/akm/sql-slog#StmtQueryContext).

### Tests

- [Test for MySQL](https://github.com/akm/sql-slog/blob/3f72cc68aefa9ac05b031d865dbdaec8a361c2c9/tests/mysql/low_level_with_context_test.go) for more details.
- [Test for PostgreSQL](https://github.com/akm/sql-slog/blob/3f72cc68aefa9ac05b031d865dbdaec8a361c2c9/tests/postgres/low_level_with_context_test.go) for more details.
- [Test for SQLite3](https://github.com/akm/sql-slog/blob/3f72cc68aefa9ac05b031d865dbdaec8a361c2c9/tests/sqlite3/low_level_without_context_test.go) for more details.

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
