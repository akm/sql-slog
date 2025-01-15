# sql-slog

![CI](https://github.com/akm/sql-slog/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/github/akm/sql-slog/graph/badge.svg?token=9BcanbSLut)](https://codecov.io/github/akm/sql-slog)
[![Go Report Card](https://goreportcard.com/badge/github.com/akm/sql-slog)](https://goreportcard.com/report/github.com/akm/sql-slog)
[![Documentation](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/akm/sql-slog)
![license](https://img.shields.io/github/license/akm/sql-slog)

A logger for Go SQL database driver without modify existing `*sql.DB` stdlib usage.

## FEATURES

- [x] Keep using (or re-use existing) `*sql.DB` as is.
- [x] No logger adapters. Just use [log/slog](https://pkg.go.dev/log/slog)
- [x] No dependencies except stdlib.
- [ ] Leveled, detailed and configurable logging.
- [ ] Duration
- [ ] Trackable log output

## INSTALL

```
go get -u github.com/akm/sql-slog
```

## USAGE

### MySQL

```golang
ctx := context.TODO()
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
db, err := sqlslog.Open(ctx, "mysql", dsn, logger)
// handle err
```

See [test for mysql](https://github.com/akm/sql-slog/blob/3f72cc68aefa9ac05b031d865dbdaec8a361c2c9/tests/mysql/low_level_with_context_test.go) for more details.

### PostgreSQL

```golang
ctx := context.TODO()
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
db, err := sqlslog.Open(ctx, "postgres", dsn, logger)
// handle err
```

See [test for postgres](https://github.com/akm/sql-slog/blob/3f72cc68aefa9ac05b031d865dbdaec8a361c2c9/tests/postgres/low_level_with_context_test.go) for more details.

### SQLite3

```golang
ctx := context.TODO()
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
db, err := sqlslog.Open(ctx, "sqlite3", dsn, logger)
// handle err
```

See [test for sqlite3](https://github.com/akm/sql-slog/blob/3f72cc68aefa9ac05b031d865dbdaec8a361c2c9/tests/sqlite3/low_level_without_context_test.go) for more details.

## MOTIVATION

I want to:

- Keep using `*sql.DB`.
  - To work with thin ORM
- Use log/slog
  - Leverage structured logging
  - Fetch and log `context.Context` value if needed

## REFERENCES

- [Stdlib sql.DB](https://github.com/golang/go/blob/master/src/database/sql/sql.go)
- [SQL driver interfaces](https://github.com/golang/go/blob/master/src/database/sql/driver/driver.go)
- [SQL driver implementation](https://go.dev/wiki/SQLDrivers)
- [log/slog](https://pkg.go.dev/log/slog)
- [Structured Logging with slog](https://go.dev/blog/slog)

## CONTRIBUTE

If you found a bug, typo, wrong test, idea, help with existing issue, or anything constructive.

Don't hesitate to create an issue or pull request.

## INSPIRED BY

- [github.com/simukti/sqldb-logger](https://github.com/simukti/sqldb-logger).

## LICENSE

[MIT](./LICENSE)
