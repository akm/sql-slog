package sqlslog

import (
	"context"
	"database/sql"
	"log/slog"
)

/*
Open opens a database specified by its driver name and a driver-specific data source name,
and returns a new database handle with logging capabilities.

ctx is the context for the open operation.
driverName is the name of the database driver, same as the driverName in [sql.Open].
dsn is the data source name, same as the dataSourceName in [sql.Open].
opts are the options for logging behavior. See [Option] for details.

The returned DB can be used the same way as *sql.DB from [sql.Open].

See the following example for usage:

[Logger]: sets the slog.Logger to be used. If not set, the default is slog.Default().

[StepOptions]: sets the options for logging behavior.

[SetStepEventMsgBuilder]: sets the function to format the step name.

[sql.Open]: https://pkg.go.dev/database/sql#Open
*/
func Open(ctx context.Context, driverName, dsn string, opts ...Option) (*sql.DB, *slog.Logger, error) { // nolint:funlen
	factory := New(driverName, dsn, opts...)
	db, err := factory.Open(ctx)
	if err != nil {
		return nil, nil, err
	}
	return db, factory.Logger(), nil
}
