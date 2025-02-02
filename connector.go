package sqlslog

import "github.com/akm/sql-slog/internal/opts"

// Set the options for Connector.Connect.
func ConnectorConnect(f func(*StepOptions)) Option {
	return opts.ConnectorConnect(f)
}
