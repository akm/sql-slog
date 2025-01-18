/*
sqlslog is a logger for Go SQL database driver without modify existing [*sql.DB] stdlib usage.
sqlslog uses [*slog.Logger] to log the SQL database driver operations.

# How to use

	ctx := context.TODO() // or a context.Context
	db, err := sqlslog.Open(ctx, "mysql", dsn)

You can also use options to customize the logger behavior.

# Options

[Logger]: sets the slog.Logger to be used. If not set, the default is slog.Default().

[StepOptions]: sets the options for the logging behavior.

[SetProcNameFormatter]: sets the function to format the step name.

[*sql.DB]: https://pkg.go.dev/database/sql#DB
[*slog.Logger]: https://pkg.go.dev/log/slog#Logger
*/
package sqlslog
