package wrap

import (
	"database/sql/driver"
	"log/slog"
)

func wrapDriver(original driver.Driver, logger *SQLLogger) driver.Driver {
	if dc, ok := original.(driver.DriverContext); ok {
		return &driverContextWrapper{
			driverWrapper: driverWrapper{original: original, logger: logger},
			original:      dc,
		}
	}
	return &driverWrapper{original: original, logger: logger}
}

// https://pkg.go.dev/database/sql/driver@go1.23.4#pkg-overview
// The driver interface has evolved over time. Drivers
// should implement Connector and DriverContext interfaces.
type driverWrapper struct {
	original driver.Driver
	logger   *SQLLogger
}

var _ driver.Driver = (*driverWrapper)(nil)

// Open implements driver.Driver.
func (w *driverWrapper) Open(dsn string) (driver.Conn, error) {
	var origConn driver.Conn
	attr, err := w.logger.With(slog.String("dsn", dsn)).StepWithoutContext(&w.logger.Options.DriverOpen, func() (*slog.Attr, error) {
		var err error
		origConn, err = w.original.Open(dsn)
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(w.logger.Options.ConnIDKey, w.logger.Options.IDGen())
		return &attrRaw, err
	})
	if err != nil {
		return nil, err
	}
	lg := w.logger
	if attr != nil {
		lg = lg.With(*attr)
	}
	return wrapConn(origConn, lg), nil
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
	attr, err := w.logger.With(slog.String("dsn", dsn)).StepWithoutContext(&w.logger.Options.DriverOpenConnector, func() (*slog.Attr, error) {
		var err error
		origConnector, err = w.original.OpenConnector(dsn)
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(w.logger.Options.ConnIDKey, w.logger.Options.IDGen())
		return &attrRaw, err
	})
	if err != nil {
		return nil, err
	}
	lg := w.logger
	if attr != nil {
		lg = lg.With(*attr)
	}
	return wrapConnector(origConnector, lg), nil
}
