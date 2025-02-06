package sqlslog

import (
	"database/sql/driver"
	"log/slog"
	"strings"
)

type driverOptions struct {
	IDGen     IDGen
	ConnIDKey string

	Open          StepOptions
	OpenConnector StepOptions

	ConnOptions      *connOptions
	ConnectorOptions *connectorOptions
}

func defaultDriverOptions(driverName string, formatter StepLogMsgFormatter) *driverOptions {
	connectorOptions := defaultConnectorOptions(driverName, formatter)
	connOptions := connectorOptions.ConnOptions
	return &driverOptions{
		IDGen:     IDGeneratorDefault,
		ConnIDKey: "conn_id",

		Open:          *defaultStepOptions(formatter, StepDriverOpen, LevelInfo),
		OpenConnector: *defaultStepOptions(formatter, StepDriverOpenConnector, LevelInfo),

		ConnOptions:      connOptions,
		ConnectorOptions: connectorOptions,
	}
}

func wrapDriver(original driver.Driver, logger *stepLogger, options *driverOptions) driver.Driver {
	driverWrapper := driverWrapper{
		original: original,
		logger:   logger,
		options:  options,
	}
	if dc, ok := original.(driver.DriverContext); ok {
		return &driverContextWrapper{
			driverWrapper: driverWrapper,
			original:      dc,
		}
	}
	return &driverWrapper
}

// https://pkg.go.dev/database/sql/driver@go1.23.4#pkg-overview
// The driver interface has evolved over time. Drivers
// should implement Connector and DriverContext interfaces.
type driverWrapper struct {
	original driver.Driver
	logger   *stepLogger
	options  *driverOptions
}

var _ driver.Driver = (*driverWrapper)(nil)

// Open implements driver.Driver.
func (w *driverWrapper) Open(dsn string) (driver.Conn, error) {
	var origConn driver.Conn
	attr, err := w.logger.With(slog.String("dsn", dsn)).StepWithoutContext(&w.options.Open, func() (*slog.Attr, error) {
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

	return wrapConn(origConn, lg, w.options.ConnOptions), nil
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
func (w *driverContextWrapper) OpenConnector(dsn string) (driver.Connector, error) { // nolint:funlen
	var origConnector driver.Connector
	attr, err := w.logger.With(slog.String("dsn", dsn)).StepWithoutContext(&w.options.OpenConnector, func() (*slog.Attr, error) {
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

	return wrapConnector(origConnector, lg, w.options.ConnectorOptions), nil
}

// DriverOpenErrorHandler returns a function that handles errors from driver.Driver.Open.
// The function returns a boolean indicating completion and a slice of slog.Attr.
//
// # For Postgres:
// If err is nil, it returns true and a slice of slog.Attr{slog.Bool("success", true)}.
// If err is io.EOF, it returns true and a slice of slog.Attr{slog.Bool("success", false)}.
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
