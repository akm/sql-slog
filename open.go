package sqlslog

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log/slog"
)

// Open opens a database specified by its database driver name and a driver-specific data source name.
// And returns a new database handle with logger.
func Open(ctx context.Context, driverName, dsn string, opts ...Option) (*sql.DB, error) {
	options := newOptions(opts...)
	logger := newLogger(options.logger, options)

	lg := logger.With(
		slog.String("driver", driverName),
		slog.String("dsn", dsn),
	)

	var db *sql.DB
	err := lg.logActionContext(ctx, &logger.options.sqlslogOpen, func() error {
		var err error
		db, err = open(driverName, dsn, logger)
		return err
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func open(driverName, dsn string, logger *logger) (*sql.DB, error) {
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
		origConnector = &dsnConnector{dsn: dsn, driver: drv}
	}

	return sql.OpenDB(wrapConnector(origConnector, logger)), nil
}
