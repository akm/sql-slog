package opts

import (
	"log/slog"
	"strings"
)

type DriverOptions struct {
	IDGen         IDGen
	ConnIDKey     string
	Open          *StepOptions
	OpenConnector *StepOptions

	Conn      *ConnOptions
	Connector *ConnectorOptions
}

const ConnIDKeyDefault = "conn_id"

func DefaultDriverOptions(driverName string, formatter StepLogMsgFormatter) *DriverOptions {
	return &DriverOptions{
		IDGen:         IDGeneratorDefault,
		ConnIDKey:     ConnIDKeyDefault,
		Open:          DefaultStepOptions(formatter, "Driver.Open", LevelInfo, DriverOpenErrorHandler(driverName)),
		OpenConnector: DefaultStepOptions(formatter, "Driver.OpenConnector", LevelInfo),

		Conn:      DefaultConnOptions(driverName, formatter),
		Connector: DefaultConnectorOptions(driverName, formatter),
	}
}

// DriverOpenErrorHandler returns a function that handles errors from driver.Driver.Open.
// The function returns a boolean indicating completion and a slice of slog.Attr.
//
// # For Postgres:
// If err is nil, it returns true and a slice of slog.Attr{slog.Bool("success", true)}.
// If err is io.EOF, it returns true and a slice of slog.Attr{slog.Bool("success", false)}.
// Otherwise, it returns false and nil.
func DriverOpenErrorHandler(driverName string) func(err error) (bool, []slog.Attr) {
	switch driverName {
	case "postgres":
		return func(err error) (bool, []slog.Attr) {
			if err == nil {
				return true, []slog.Attr{slog.Bool("success", true)}
			}
			if strings.ToUpper(err.Error()) == "EOF" {
				return true, []slog.Attr{slog.Bool("success", false)}
			}
			return false, nil
		}
	default:
		return nil
	}
}
