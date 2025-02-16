package sqlslog

import (
	"context"
	"database/sql"
	"log/slog"
)

type Factory struct {
	options    *options
	driverName string
	dsn        string
}

func New(driverName, dsn string, opts ...Option) *Factory {
	options := newOptions(driverName, opts...)
	return &Factory{driverName: driverName, dsn: dsn, options: options}
}

func (f *Factory) Logger() *slog.Logger {
	return f.options.logger
}

func (f *Factory) Open(ctx context.Context) (*sql.DB, error) {
	logger := newStepLogger(&stepLoggerOptions{
		logger:       f.options.logger,
		durationKey:  f.options.durationKey,
		durationType: f.options.durationType,
	})
	return open(ctx, f.driverName, f.dsn, logger, &f.options.sqlslogOptions)
}
