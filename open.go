package sqlslog

import (
	"context"
	"database/sql"

	"github.com/akm/sql-slog/internal/opts"
	"github.com/akm/sql-slog/internal/wrap"
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
func Open(ctx context.Context, driverName, dsn string, options ...Option) (*sql.DB, error) {
	gOptions := opts.NewOptions(driverName, options...)
	stepLogger := wrap.NewStepLogger(gOptions.Logger, wrap.DurationAttrFunc(gOptions.DurationKey, gOptions.DurationType))

	openOptions := buildOpenOptions(gOptions)

	return wrap.Open(ctx, driverName, dsn, stepLogger, openOptions)
}

func buildOpenOptions(options *Options) *wrap.OpenOptions {
	connOptions := &wrap.ConnOptions{
		IDGen:   options.IDGen,
		Begin:   &options.ConnBegin,
		BeginTx: &options.ConnBeginTx,
		TxIDKey: options.TxIDKey,
		Tx: &wrap.TxOptions{
			Commit:   &options.TxCommit,
			Rollback: &options.TxRollback,
		},
		Close:          &options.ConnClose,
		Prepare:        &options.ConnPrepare,
		PrepareContext: &options.ConnPrepareContext,
		StmtIDKey:      options.StmtIDKey,
		Stmt: &wrap.StmtOptions{
			Close:        &options.StmtClose,
			Exec:         &options.StmtExec,
			Query:        &options.StmtQuery,
			ExecContext:  &options.StmtExecContext,
			QueryContext: &options.StmtQueryContext,
			Rows: &wrap.RowsOptions{
				Close:         &options.RowsClose,
				Next:          &options.RowsNext,
				NextResultSet: &options.RowsNextResultSet,
			},
		},
		ResetSession: &options.ConnResetSession,
		Ping:         &options.ConnPing,
		ExecContext:  &options.ConnExecContext,
		QueryContext: &options.ConnQueryContext,
		Rows: &wrap.RowsOptions{
			Close:         &options.RowsClose,
			Next:          &options.RowsNext,
			NextResultSet: &options.RowsNextResultSet,
		},
	}
	return &wrap.OpenOptions{
		Open: &options.SqlslogOpen,
		Driver: &wrap.DriverOptions{
			IDGen:         options.IDGen,
			ConnIDKey:     options.ConnIDKey,
			Open:          &options.DriverOpen,
			OpenConnector: &options.DriverOpenConnector,
			Conn:          connOptions,
			Connector: &wrap.ConnectorOptions{
				Connect: &options.ConnectorConnect,
				Conn:    connOptions,
			},
		},
	}
}

// Set the options for sqlslog.Open.
func SqlslogOpen(f func(*StepOptions)) Option { return opts.SqlslogOpen(f) } // nolint:revive
