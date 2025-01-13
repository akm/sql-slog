package sqlslog

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log/slog"
)

func Open(ctx context.Context, driverName, dsn string, logger *slog.Logger) (*sql.DB, error) {
	lg := logger.With(
		slog.String("driver", driverName),
		slog.String("dsn", dsn),
	)
	lg.Debug("sqlslog.Open Start")

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		lg.ErrorContext(ctx, "sqlslog.Open Error", "error", err)
		return nil, err
	}
	// This db is not used directly, but it is used to get the driver.

	drv := wrapDriver(db.Driver(), logger)

	var origConnector driver.Connector

	if dc, ok := drv.(driver.DriverContext); ok {
		connector, err := dc.OpenConnector(dsn)
		if err != nil {
			lg.ErrorContext(ctx, "sqlslog.Open OpenConnector Error", "error", err)
			return nil, err
		}
		origConnector = connector
	} else {
		origConnector = &dsnConnector{dsn: dsn, driver: drv}
	}

	lg.Debug("sqlslog.Open Complete")
	return sql.OpenDB(wrapConnector(origConnector, logger)), nil
}
