/*
sqlslog is a logger for Go SQL database drivers without modifying existing [*sql.DB] stdlib usage.
sqlslog uses [*slog.Logger] to log SQL database driver operations.

# How to use

	ctx := context.TODO() // or a context.Context
	db, err := sqlslog.Open(ctx, "mysql", dsn)

You can also use options to customize the logger's behavior.

[Open] takes [opts.Option] s to customize the logging behavior.
[opts.Option] is created by using functions like [opts.Logger], [opts.ConnPrepareContext], [opts.StmtQueryContext], etc.

# Logger

[opts.Logger] sets the slog.Logger to be used. If not set, the default is slog.Default().

The logger from slog.Default() does not have options for the log levels [opts.LevelTrace] and [opts.LevelVerbose].

You can create a [slog.Handler] by using [opts.NewTextHandler] or [opts.NewJSONHandler] customized for sqlslog [opts.Level].
See examples for [opts.NewTextHandler] and [opts.NewJSONHandler] for more details.

# Step

In sqlslog terms, a step is a logical operation in the database driver, such as a query, a ping, a prepare, etc.

A step has three events: start, error, and complete.

sqlslog provides a way to customize the log message and log level for each step event.

You can customize them by using functions that take [opts.StepOptions] and return [opts.Option], like [opts.ConnPrepareContext] or [opts.StmtQueryContext].

# Default Step Log Message Formatter

The default step log message formatter is [opts.StepLogMsgWithEventName].

You can change the default step log message formatter by calling [opts.SetStepLogMsgFormatter].

[*sql.DB]: https://pkg.go.dev/database/sql#DB
[*slog.Logger]: https://pkg.go.dev/log/slog#Logger
[slog.Handler]: https://pkg.go.dev/log/slog#Handler
*/
package sqlslog
