package sqlslog

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log/slog"
)

// Open opens a database specified by its database driver name and a driver-specific data source name.
// And returns a new database handle with logger.
func Open(ctx context.Context, driverName, dsn string, logger *slog.Logger) (*sql.DB, error) {
	lg := logger.With(
		slog.String("driver", driverName),
		slog.String("dsn", dsn),
	)
	lg.DebugContext(ctx, "sqlslog.Open Start")

	db, err := open(driverName, dsn, logger)
	if err != nil {
		lg.ErrorContext(ctx, "sqlslog.Open Error", "error", err)
		return nil, err
	}

	lg.InfoContext(ctx, "sqlslog.Open Complete")
	return db, nil
}

func open(driverName, dsn string, logger *slog.Logger) (*sql.DB, error) {
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
