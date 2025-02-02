package wrap

import (
	"database/sql/driver"
	"log/slog"

	"github.com/akm/sql-slog/internal/opts"
)

type DriverOptions = opts.DriverOptions

func WrapDriver(original driver.Driver, logger *StepLogger, options *DriverOptions) driver.Driver { // nolint:revive
	r := driverWrapper{
		original: original,
		logger:   logger,
		options:  options,
	}
	if dc, ok := original.(driver.DriverContext); ok {
		return &driverContextWrapper{
			driverWrapper: r,
			original:      dc,
		}
	}
	return &r
}

// https://pkg.go.dev/database/sql/driver@go1.23.4#pkg-overview
// The driver interface has evolved over time. Drivers
// should implement Connector and DriverContext interfaces.
type driverWrapper struct {
	original driver.Driver
	logger   *StepLogger
	options  *DriverOptions
}

var _ driver.Driver = (*driverWrapper)(nil)

// Open implements driver.Driver.
func (w *driverWrapper) Open(dsn string) (driver.Conn, error) {
	var origConn driver.Conn
	attr, err := w.logger.With(slog.String("dsn", dsn)).StepWithoutContext(w.options.Open, func() (*slog.Attr, error) {
		var err error
		origConn, err = w.original.Open(dsn)
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(w.options.ConnIDKey, w.options.IDGen())
		return &attrRaw, err
	})
	if err != nil {
		return nil, err
	}
	lg := w.logger
	if attr != nil {
		lg = lg.With(*attr)
	}
	return WrapConn(origConn, lg, w.options.Conn), nil
}

type driverContextWrapper struct {
	driverWrapper
	original driver.DriverContext
}

var (
	_ driver.Driver        = (*driverContextWrapper)(nil)
	_ driver.DriverContext = (*driverContextWrapper)(nil)
)

// var _ driver.Connector = (*driverWrapper)(nil)

// OpenConnector implements driver.DriverContext.
func (w *driverContextWrapper) OpenConnector(dsn string) (driver.Connector, error) {
	var origConnector driver.Connector
	attr, err := w.logger.With(slog.String("dsn", dsn)).StepWithoutContext(w.options.OpenConnector, func() (*slog.Attr, error) {
		var err error
		origConnector, err = w.original.OpenConnector(dsn)
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(w.options.ConnIDKey, w.options.IDGen())
		return &attrRaw, err
	})
	if err != nil {
		return nil, err
	}
	lg := w.logger
	if attr != nil {
		lg = lg.With(*attr)
	}
	return WrapConnector(origConnector, lg, w.options.Connector), nil
}
