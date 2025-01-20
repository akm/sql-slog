package sqlslog

import (
	"database/sql/driver"
	"errors"
	"io"
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
	lg := w.logger.With(slog.String("dsn", dsn))
	var origConn driver.Conn
	err := lg.StepWithoutContext(&w.logger.options.driverOpen, func() error {
		var err error
		origConn, err = w.original.Open(dsn)
		return err
	})
	if err != nil {
		return nil, err
	}
	return wrapConn(origConn, w.logger), nil
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
	lg := w.logger.With(slog.String("dsn", dsn))
	var origConnector driver.Connector
	err := lg.StepWithoutContext(&w.logger.options.driverOpenConnector, func() error {
		var err error
		origConnector, err = w.original.OpenConnector(dsn)
		return err
	})
	if err != nil {
		return nil, err
	}
	return wrapConnector(origConnector, w.logger), nil
}

// HandleDriverOpenError returns completed and slice of slog.Attr.
// If err is nil, it returns true and a slice of slog.Attr{slog.Bool("eof", false)}.
// If err is io.EOF, it returns true and a slice of slog.Attr{slog.Bool("eof", true)}.
// Otherwise, it returns false and nil.
func HandleDriverOpenError(err error) (bool, []slog.Attr) {
	if err == nil {
		return true, []slog.Attr{slog.Bool("eof", false)}
	}
	if errors.Is(err, io.EOF) {
		return true, []slog.Attr{slog.Bool("eof", true)}
	} else if strings.ToUpper(err.Error()) == "EOF" {
		return true, []slog.Attr{slog.Bool("eof", true)}
	}
	return false, nil
}
