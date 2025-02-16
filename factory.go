package sqlslog

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log/slog"
)

type factoryOptions struct {
	Open          StepOptions
	DriverOptions *driverOptions
	SlogOptions   *slogOptions
}

func defaultFactoryOptions(driverName string, msgb StepEventMsgBuilder) *factoryOptions {
	driverOptions := defaultDriverOptions(driverName, msgb)
	return &factoryOptions{
		Open:          *defaultStepOptions(msgb, StepSqlslogOpen, LevelInfo),
		DriverOptions: driverOptions,
		SlogOptions:   defaultSlogOptions(),
	}
}

type Factory struct {
	options    *options
	driverName string
	dsn        string

	handler slog.Handler
	logger  *slog.Logger
}

func New(driverName, dsn string, opts ...Option) *Factory {
	options := newOptions(driverName, opts...)
	return &Factory{driverName: driverName, dsn: dsn, options: options}
}

func (f *Factory) Handler() slog.Handler {
	if f.handler == nil {
		f.handler = f.options.factoryOptions.SlogOptions.newHandler()
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
	logger := newStepLogger(&stepLoggerOptions{
		logger:       f.options.logger,
		durationKey:  f.options.durationKey,
		durationType: f.options.durationType,
	})
	return open(ctx, f.driverName, f.dsn, logger, &f.options.factoryOptions)
}

func open(ctx context.Context, driverName, dsn string, logger *stepLogger, options *factoryOptions) (*sql.DB, error) {
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
