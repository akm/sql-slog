/*
sqlslog is a logger for Go SQL database drivers without modifying existing [*sql.DB] stdlib usage.
sqlslog uses [*slog.Logger] to log SQL database driver operations.

# How to use

	ctx := context.TODO() // or a context.Context
	db, err := sqlslog.Open(ctx, "mysql", dsn)

You can also use options to customize the logger's behavior.

[Open] takes [Option] s to customize the logging behavior.
[Option] is created by using functions like [Logger], [ConnPrepareContext], [StmtQueryContext], etc.

# Logger

[Logger] sets the slog.Logger to be used. If not set, the default is slog.Default().

The logger from slog.Default() does not have options for the log levels [LevelTrace] and [LevelVerbose].

You can create a [slog.Handler] by using [sqlslog.NewTextHandler] or [sqlslog.NewJSONHandler] customized for sqlslog [Level].
See examples for [NewTextHandler] and [NewJSONHandler] for more details.

# Level

sqlslog has 6 log levels: [LevelVerbose], [LevelTrace], [LevelDebug], [LevelInfo], [LevelWarn], and [LevelError].
[LevelDebug], [LevelInfo], [LevelWarn], and [LevelError] are the same as slog's log levels.
[LevelVerbose] and [LevelTrace] are extra log levels for sqlslog.
[LevelVerbose] is the lowest log level, and [LevelTrace] is the second lowest log level.

# Step

In sqlslog terms, a step is a logical operation in the database driver, such as a query, a ping, a prepare, etc.

A step has three events: start, error, and complete.

sqlslog provides a way to customize the log message and log level for each step event.

You can customize them by using functions that take [StepOptions] and return [Option], like [ConnPrepareContext] or [StmtQueryContext].

# Default Step Log Message Formatter

The default step log message formatter is [StepLogMsgWithEventName].

You can change the default step log message formatter by calling [SetStepLogMsgFormatter].

[*sql.DB]: https://pkg.go.dev/database/sql#DB
[*slog.Logger]: https://pkg.go.dev/log/slog#Logger
[slog.Handler]: https://pkg.go.dev/log/slog#Handler
*/
package sqlslog
