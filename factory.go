package sqlslog

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log/slog"
)

type Factory struct {
	options    *options
	driverName string
	dsn        string

	handler slog.Handler
	logger  *slog.Logger
}

func New(driverName, dsn string, opts ...Option) *Factory {
	options := newOptions(driverName, opts...)
	return &Factory{
		driverName: driverName,
		dsn:        dsn,
		options:    options,
		handler:    options.SlogOptions.handler,
	}
}

func (f *Factory) Handler() slog.Handler {
	if f.handler == nil {
		f.handler = f.options.SlogOptions.newHandler()
	}
	return f.handler
}

func (f *Factory) Logger() *slog.Logger {
	if f.logger == nil {
		f.logger = slog.New(f.Handler())
	}
	return f.logger
}

func (f *Factory) Open(ctx context.Context) (*sql.DB, error) {
	stepLogger := newStepLogger(f.Logger(), &f.options.stepLoggerOptions)
	return open(ctx, f.driverName, f.dsn, stepLogger, f.options)
}

func open(ctx context.Context, driverName, dsn string, logger *stepLogger, options *options) (*sql.DB, error) {
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
