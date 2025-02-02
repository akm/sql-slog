package wrap

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log/slog"

	"github.com/akm/sql-slog/internal/opts"
)

type OpenOptions = opts.OpenOptions

var DefaultOpenOptions = opts.DefaultOpenOptions

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
	options := NewOptions(driverName, opts...)
	logger := NewStepLogger(options.Logger, DurationAttrFunc(options.DurationKey, options.DurationType))

	openOptions := buildOpenOptions(options)

	lg := logger.With(
		slog.String("driver", driverName),
		slog.String("dsn", dsn),
	)

	var db *sql.DB
	err := IgnoreAttr(lg.Step(ctx, openOptions.Open, func() (*slog.Attr, error) {
		var err error
		db, err = open(driverName, dsn, logger, openOptions)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func open(driverName, dsn string, logger *StepLogger, options *OpenOptions) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	// This db is not used directly, but it is used to get the driver.
	drv := WrapDriver(db.Driver(), logger, options.Driver)
	return openWithDriver(drv, dsn, logger, options.Driver.Connector)
}

func openWithDriver(drv driver.Driver, dsn string, logger *StepLogger, connectorOptions *ConnectorOptions) (*sql.DB, error) {
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
	return sql.OpenDB(WrapConnector(origConnector, logger, connectorOptions)), nil
}

func buildOpenOptions(options *Options) *OpenOptions {
	connOptions := &ConnOptions{
		IDGen:   options.IDGen,
		Begin:   &options.ConnBegin,
		BeginTx: &options.ConnBeginTx,
		TxIDKey: options.TxIDKey,
		Tx: &TxOptions{
			Commit:   &options.TxCommit,
			Rollback: &options.TxRollback,
		},
		Close:          &options.ConnClose,
		Prepare:        &options.ConnPrepare,
		PrepareContext: &options.ConnPrepareContext,
		StmtIDKey:      options.StmtIDKey,
		Stmt: &StmtOptions{
			Close:        &options.StmtClose,
			Exec:         &options.StmtExec,
			Query:        &options.StmtQuery,
			ExecContext:  &options.StmtExecContext,
			QueryContext: &options.StmtQueryContext,
			Rows: &RowsOptions{
				Close:         &options.RowsClose,
				Next:          &options.RowsNext,
				NextResultSet: &options.RowsNextResultSet,
			},
		},
		ResetSession: &options.ConnResetSession,
		Ping:         &options.ConnPing,
		ExecContext:  &options.ConnExecContext,
		QueryContext: &options.ConnQueryContext,
		Rows: &RowsOptions{
			Close:         &options.RowsClose,
			Next:          &options.RowsNext,
			NextResultSet: &options.RowsNextResultSet,
		},
	}
	return &OpenOptions{
		Open: &options.SqlslogOpen,
		Driver: &DriverOptions{
			IDGen:         options.IDGen,
			ConnIDKey:     options.ConnIDKey,
			Open:          &options.DriverOpen,
			OpenConnector: &options.DriverOpenConnector,
			Conn:          connOptions,
			Connector: &ConnectorOptions{
				Connect: &options.ConnectorConnect,
				Conn:    connOptions,
			},
		},
	}
}
