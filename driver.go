package sqlslog

import (
	"database/sql/driver"
	"log/slog"
	"strings"
)

func wrapDriver(original driver.Driver, logger *stepLogger) driver.Driver {
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
	logger   *stepLogger
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

	txOptions := &txOptions{
		Commit:   w.logger.options.txCommit,
		Rollback: w.logger.options.txRollback,
	}
	rowOptions := &rowsOptions{
		Close:         w.logger.options.rowsClose,
		Next:          w.logger.options.rowsNext,
		NextResultSet: w.logger.options.rowsNextResultSet,
	}
	stmtOptions := &stmtOptions{
		Close:        w.logger.options.stmtClose,
		Exec:         w.logger.options.stmtExec,
		Query:        w.logger.options.stmtQuery,
		ExecContext:  w.logger.options.stmtExecContext,
		QueryContext: w.logger.options.stmtQueryContext,
		Rows:         rowOptions,
	}
	connOptions := &connOptions{
		IDGen: w.logger.options.idGen,

		Begin:     w.logger.options.connBegin,
		BeginTx:   w.logger.options.connBeginTx,
		TxIDKey:   w.logger.options.txIDKey,
		TxOptions: txOptions,

		Close: w.logger.options.connClose,

		Prepare:        w.logger.options.connPrepare,
		PrepareContext: w.logger.options.connPrepareContext,
		StmtIDKey:      w.logger.options.stmtIDKey,
		StmtOptions:    stmtOptions,

		ResetSession: w.logger.options.connResetSession,
		Ping:         w.logger.options.connPing,

		ExecContext: w.logger.options.connExecContext,

		QueryContext: w.logger.options.connQueryContext,
		RowsOptions:  rowOptions,
	}
	return wrapConn(origConn, lg, connOptions), nil
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
