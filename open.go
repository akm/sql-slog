package sqlslog

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log/slog"
)

/*
Open opens a database specified by its driver name and a driver-specific data source name,
and returns a new database handle with logging capabilities.

ctx is the context for the open operation.
driverName is the name of the database driver, same as the driverName in [sql.Open].
dsn is the data source name, same as the dataSourceName in [sql.Open].
opts are the options for logging behavior. See [Option] for details.

The returned DB can be used the same way as *sql.DB from [sql.Open].

See the following example for usage:

[Logger]: sets the slog.Logger to be used. If not set, the default is slog.Default().

[StepOptions]: sets the options for logging behavior.

[SetStepLogMsgFormatter]: sets the function to format the step name.

[sql.Open]: https://pkg.go.dev/database/sql#Open
*/
func Open(ctx context.Context, driverName, dsn string, opts ...Option) (*sql.DB, error) { // nolint:funlen
	options := newOptions(driverName, opts...)
	logger := newStepLogger(options.logger, options)

	txOptions := &txOptions{
		Commit:   options.txCommit,
		Rollback: options.txRollback,
	}
	rowOptions := &rowsOptions{
		Close:         options.rowsClose,
		Next:          options.rowsNext,
		NextResultSet: options.rowsNextResultSet,
	}
	stmtOptions := &stmtOptions{
		Close:        options.stmtClose,
		Exec:         options.stmtExec,
		Query:        options.stmtQuery,
		ExecContext:  options.stmtExecContext,
		QueryContext: options.stmtQueryContext,
		Rows:         rowOptions,
	}
	connOptions := &connOptions{
		IDGen: options.idGen,

		Begin:     options.connBegin,
		BeginTx:   options.connBeginTx,
		TxIDKey:   options.txIDKey,
		TxOptions: txOptions,

		Close: options.connClose,

		Prepare:        options.connPrepare,
		PrepareContext: options.connPrepareContext,
		StmtIDKey:      options.stmtIDKey,
		StmtOptions:    stmtOptions,

		ResetSession: options.connResetSession,
		Ping:         options.connPing,

		ExecContext: options.connExecContext,

		QueryContext: options.connQueryContext,
		RowsOptions:  rowOptions,
	}
	connectorOptions := &connectorOptions{
		Connect:     options.connectorConnect,
		ConnOptions: connOptions,
	}
	driverOptions := &driverOptions{
		IDGen:            options.idGen,
		ConnIDKey:        options.connIDKey,
		Open:             options.driverOpen,
		OpenConnector:    options.driverOpenConnector,
		ConnOptions:      connOptions,
		ConnectorOptions: connectorOptions,
	}

	lg := logger.With(
		slog.String("driver", driverName),
		slog.String("dsn", dsn),
	)

	var db *sql.DB
	err := ignoreAttr(lg.Step(ctx, &logger.options.sqlslogOpen, func() (*slog.Attr, error) {
		var err error
		db, err = openWithDriver(driverName, dsn, logger, driverOptions)
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
