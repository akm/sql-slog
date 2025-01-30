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
func Open(ctx context.Context, driverName, dsn string, opts ...Option) (*sql.DB, error) {
	options := newOptions(driverName, opts...)
	logger := newLogger(options.logger, options)

	lg := logger.With(
		slog.String("driver", driverName),
		slog.String("dsn", dsn),
	)

	var db *sql.DB
	err := ignoreAttr(lg.Step(ctx, &logger.options.sqlslogOpen, func() (*slog.Attr, error) {
		var err error
		db, err = open(driverName, dsn, logger)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func open(driverName, dsn string, logger *logger) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	// This db is not used directly, but it is used to get the driver.

	drv := wrapDriver(db.Driver(), logger)

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

	opts := logger.options
	options := &connectorOptions{
		Connect: &opts.connectorConnect,
		Conn: &connOptions{
			idGen:   opts.idGen,
			Begin:   &opts.connBegin,
			BeginTx: &opts.connBeginTx,
			txIDKey: opts.txIDKey,
			Tx: &txOptions{
				Commit:   &opts.txCommit,
				Rollback: &opts.txRollback,
			},
			Close:          &opts.connClose,
			Prepare:        &opts.connPrepare,
			PrepareContext: &opts.connPrepareContext,
			stmtIDKey:      opts.stmtIDKey,
			Stmt: &stmtOptions{
				Close:        &opts.stmtClose,
				Exec:         &opts.stmtExec,
				Query:        &opts.stmtQuery,
				ExecContext:  &opts.stmtExecContext,
				QueryContext: &opts.stmtQueryContext,
				Rows: &rowsOptions{
					Close:         &opts.rowsClose,
					Next:          &opts.rowsNext,
					NextResultSet: &opts.rowsNextResultSet,
				},
			},
			ResetSession: &opts.connResetSession,
			Ping:         &opts.connPing,
			ExecContext:  &opts.connExecContext,
			QueryContext: &opts.connQueryContext,
			Rows: &rowsOptions{
				Close:         &opts.rowsClose,
				Next:          &opts.rowsNext,
				NextResultSet: &opts.rowsNextResultSet,
			},
		},
	}
	return sql.OpenDB(wrapConnector(origConnector, logger, options)), nil
}
