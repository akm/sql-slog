package sqlslog

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log/slog"
)

type openOptions struct {
	Open   *StepOptions
	Driver *DriverOptions
}

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
	logger := newLogger(options.logger, DurationAttrFunc(options.durationKey, options.durationType))

	openOptions := buildOpenOptions(options)

	lg := logger.With(
		slog.String("driver", driverName),
		slog.String("dsn", dsn),
	)

	var db *sql.DB
	err := ignoreAttr(lg.Step(ctx, openOptions.Open, func() (*slog.Attr, error) {
		var err error
		db, err = open(driverName, dsn, logger, openOptions)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func open(driverName, dsn string, logger *logger, options *openOptions) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	// This db is not used directly, but it is used to get the driver.
	drv := WrapDriver(db.Driver(), logger, options.Driver)
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

	return sql.OpenDB(WrapConnector(origConnector, logger, options.Driver.Connector)), nil
}

func buildOpenOptions(options *options) *openOptions {
	connOptions := &ConnOptions{
		idGen:   options.idGen,
		Begin:   &options.connBegin,
		BeginTx: &options.connBeginTx,
		txIDKey: options.txIDKey,
		Tx: &TxOptions{
			Commit:   &options.txCommit,
			Rollback: &options.txRollback,
		},
		Close:          &options.connClose,
		Prepare:        &options.connPrepare,
		PrepareContext: &options.connPrepareContext,
		stmtIDKey:      options.stmtIDKey,
		Stmt: &StmtOptions{
			Close:        &options.stmtClose,
			Exec:         &options.stmtExec,
			Query:        &options.stmtQuery,
			ExecContext:  &options.stmtExecContext,
			QueryContext: &options.stmtQueryContext,
			Rows: &RowsOptions{
				Close:         &options.rowsClose,
				Next:          &options.rowsNext,
				NextResultSet: &options.rowsNextResultSet,
			},
		},
		ResetSession: &options.connResetSession,
		Ping:         &options.connPing,
		ExecContext:  &options.connExecContext,
		QueryContext: &options.connQueryContext,
		Rows: &RowsOptions{
			Close:         &options.rowsClose,
			Next:          &options.rowsNext,
			NextResultSet: &options.rowsNextResultSet,
		},
	}
	return &openOptions{
		Open: &options.sqlslogOpen,
		Driver: &DriverOptions{
			IDGen:         options.idGen,
			connIDKey:     options.connIDKey,
			Open:          &options.driverOpen,
			OpenConnector: &options.driverOpenConnector,
			Conn:          connOptions,
			Connector: &ConnectorOptions{
				Connect: &options.connectorConnect,
				Conn:    connOptions,
			},
		},
	}
}
