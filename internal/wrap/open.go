package wrap

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log/slog"

	"github.com/akm/sql-slog/sqlslogopts"
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

[SetStepLogMsgFormatter]: sets the function to format the step name.

[sql.Open]: https://pkg.go.dev/database/sql#Open
*/
func Open(ctx context.Context, driverName, dsn string, opts ...sqlslogopts.Option) (*sql.DB, error) {
	options := NewOptions(driverName, opts...)
	logger := NewSQLLogger(options.Logger, options)

	lg := logger.With(
		slog.String("driver", driverName),
		slog.String("dsn", dsn),
	)

	var db *sql.DB
	err := IgnoreAttr(lg.Step(ctx, &logger.Options.SqlslogOpen, func() (*slog.Attr, error) {
		var err error
		db, err = openWithWrappedConnector(driverName, dsn, logger)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func openWithWrappedConnector(driverName, dsn string, logger *SQLLogger) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	// This db is not used directly, but it is used to get the driver.

	drv := wrapDriver(db.Driver(), logger)

	var origConnector driver.Connector

	if dc, ok := drv.(driver.DriverContext); ok {
		connector, err := dc.OpenConnector(dsn)
		if err != nil {
			return nil, err
		}
		origConnector = connector
	} else {
		origConnector = NewDsnConnector(dsn, drv)
	}

	return sql.OpenDB(wrapConnector(origConnector, logger)), nil
}
