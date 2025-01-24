package sqlslog

import (
	"database/sql/driver"
	"log/slog"
	"strings"
)

func wrapDriver(original driver.Driver, logger *logger) driver.Driver {
	if dc, ok := original.(driver.DriverContext); ok {
		return &driverContextWrapper{
			driverWrapper: driverWrapper{original: original, logger: logger},
			original:      dc,
		}
	} else {
		return &driverWrapper{original: original, logger: logger}
	}
}

// https://pkg.go.dev/database/sql/driver@go1.23.4#pkg-overview
// The driver interface has evolved over time. Drivers
// should implement Connector and DriverContext interfaces.
type driverWrapper struct {
	original driver.Driver
	logger   *logger
}

var _ driver.Driver = (*driverWrapper)(nil)

// Open implements driver.Driver.
func (w *driverWrapper) Open(dsn string) (driver.Conn, error) {
	var origConn driver.Conn
	attr, err := w.logger.With(slog.String("dsn", dsn)).StepWithoutContext(&w.logger.options.driverOpen, func() (*slog.Attr, error) {
		var err error
		origConn, err = w.original.Open(dsn)
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(w.logger.options.connIDKey, w.logger.options.idGen())
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

var _ driver.Driver = (*driverContextWrapper)(nil)
var _ driver.DriverContext = (*driverContextWrapper)(nil)

// var _ driver.Connector = (*driverWrapper)(nil)

// OpenConnector implements driver.DriverContext.
func (w *driverContextWrapper) OpenConnector(dsn string) (driver.Connector, error) {
	var origConnector driver.Connector
	attr, err := w.logger.With(slog.String("dsn", dsn)).StepWithoutContext(&w.logger.options.driverOpenConnector, func() (*slog.Attr, error) {
		var err error
		origConnector, err = w.original.OpenConnector(dsn)
		if err != nil {
			return nil, err
		}
		attrRaw := slog.String(w.logger.options.connIDKey, w.logger.options.idGen())
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

// DriverOpenErrorHandler returns a function that handles the error of driver.Driver.Open.
// The function returns completed and slice of slog.Attr.
//
// # For Postgres:
// If err is nil, it returns true and a slice of slog.Attr{slog.Bool("eof", false)}.
// If err is io.EOF, it returns true and a slice of slog.Attr{slog.Bool("eof", true)}.
// Otherwise, it returns false and nil.
func DriverOpenErrorHandler(driverName string) func(err error) (bool, []slog.Attr) {
	switch driverName {
	case "postgres":
		return func(err error) (bool, []slog.Attr) {
			if err == nil {
				return true, []slog.Attr{slog.Bool("success", true)}
			}
			if strings.ToUpper(err.Error()) == "EOF" {
				return true, []slog.Attr{slog.Bool("success", false)}
			}
			return false, nil
		}
	default:
		return nil
	}
}
