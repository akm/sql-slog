package sqlslog

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log/slog"
)

type sqlslogOptions struct {
	Open          StepOptions
	DriverOptions *driverOptions
}

func defaultSqlslogOptions(driverName string, formatter StepLogMsgFormatter) *sqlslogOptions {
	driverOptions := defaultDriverOptions(driverName, formatter)
	return &sqlslogOptions{
		Open:          *defaultStepOptions(formatter, StepSqlslogOpen, LevelInfo),
		DriverOptions: driverOptions,
	}
}

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
func Open(ctx context.Context, driverName, dsn string, opts ...Option) (*sql.DB, error) { // nolint:funlen
	options := newOptions(driverName, opts...)
	logger := newStepLogger(&stepLoggerOptions{
		logger:       options.logger,
		durationKey:  options.durationKey,
		durationType: options.durationType,
	})

	return open(ctx, driverName, dsn, logger, &options.sqlslogOptions)
}

func open(ctx context.Context, driverName, dsn string, logger *stepLogger, options *sqlslogOptions) (*sql.DB, error) {
	lg := logger.With(
		slog.String("driver", driverName),
		slog.String("dsn", dsn),
	)

	var db *sql.DB
	err := ignoreAttr(lg.Step(ctx, &options.Open, func() (*slog.Attr, error) {
		var err error
		db, err = openWithDriver(driverName, dsn, logger, options.DriverOptions)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func openWithDriver(driverName, dsn string, logger *stepLogger, driverOptions *driverOptions) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	// This db is not used directly, but it is used to get the driver.

	drv := wrapDriver(db.Driver(), logger, driverOptions)

	return openWithWrappedDriver(drv, dsn, logger, driverOptions)
}

func openWithWrappedDriver(drv driver.Driver, dsn string, logger *stepLogger, driverOptions *driverOptions) (*sql.DB, error) {
	var origConnector driver.Connector

	if dc, ok := drv.(driver.DriverContext); ok {
		connector, err := dc.OpenConnector(dsn)
		if err != nil {
			return nil, err
		}
		origConnector = connector
	} else {
		origConnector = &dsnConnector{dsn: dsn, driver: drv}
	}

	return sql.OpenDB(wrapConnector(origConnector, logger, driverOptions.ConnectorOptions)), nil
}
