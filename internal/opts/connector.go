package opts

import (
	"errors"
	"io"
	"log/slog"
)

type ConnectorOptions struct {
	Connect *StepOptions

	Conn *ConnOptions
}

func DefaultConnectorOptions(driverName string, formatter StepLogMsgFormatter) *ConnectorOptions {
	return &ConnectorOptions{
		Connect: DefaultStepOptions(formatter, "Connector.Connect", LevelInfo, ConnectorConnectErrorHandler(driverName)),
		Conn:    DefaultConnOptions(driverName, formatter),
	}
}

// ConnectorConnectErrorHandler returns a function that handles errors from driver.Connector.Connect.
// The function returns a boolean indicating completion and a slice of slog.Attr.
//
// # For Postgres:
// If err is nil, it returns true and a slice of slog.Attr{slog.Bool("success", true)}.
// If err is io.EOF, it returns true and a slice of slog.Attr{slog.Bool("success", false)}.
// Otherwise, it returns false and nil.
func ConnectorConnectErrorHandler(driverName string) func(err error) (bool, []slog.Attr) {
	switch driverName {
	case "mysql":
		return func(err error) (bool, []slog.Attr) {
			if err == nil {
				return true, []slog.Attr{slog.Bool("success", true)}
			}
			if err.Error() == "driver: bad connection" {
				return true, []slog.Attr{slog.Bool("success", false)}
			}
			return false, nil
		}
	case "postgres":
		return func(err error) (bool, []slog.Attr) {
			if err == nil {
				return true, []slog.Attr{slog.Bool("success", true)}
			}
			if errors.Is(err, io.EOF) {
				return true, []slog.Attr{slog.Bool("success", false)}
			}
			return false, nil
		}
	default:
		return nil
	}
}

// Set the options for Connector.Connect.
func ConnectorConnect(f func(*StepOptions)) Option {
	return func(o *Options) { f(o.Driver.Connector.Connect) }
}
