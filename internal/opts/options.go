package opts

import (
	"log/slog"
)

type Options struct {
	Logger *slog.Logger

	DurationKey  string
	DurationType DurationType

	OpenOptions
}

func NewDefaultOptions(driverName string, formatter StepLogMsgFormatter) *Options {
	openOptions := DefaultOpenOptions(driverName, formatter)
	// driverOptions := openOptions.Driver
	// connectorOptions := driverOptions.Connector
	// connOptions := connectorOptions.Conn
	// stmtOptions := connOptions.Stmt
	// rowsOptions := stmtOptions.Rows
	// txOptions := connOptions.Tx

	return &Options{
		Logger:       slog.Default(),
		DurationKey:  DurationKeyDefault,
		DurationType: DurationNanoSeconds,

		OpenOptions: *openOptions,
	}
}

// DurationKeyDefault is the default key for duration value in log.
const DurationKeyDefault = "duration"

// Option is a function that sets an option on the options struct.
type Option func(*Options)

var stepLogMsgFormatter = StepLogMsgWithoutEventName

// SetStepLogMsgFormatter sets the formatter for the step name used in logs.
// If not set, the default is StepLogMsgWithEventName.
func SetStepLogMsgFormatter(f StepLogMsgFormatter) { stepLogMsgFormatter = f }

func NewOptions(driverName string, opts ...Option) *Options {
	o := NewDefaultOptions(driverName, stepLogMsgFormatter)
	for _, opt := range opts {
		opt(o)
	}
	return o
}
