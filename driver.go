package sqlslog

import "github.com/akm/sql-slog/internal/opts"

// Set the options for Driver.Open.
func DriverOpen(f func(*StepOptions)) Option { return opts.DriverOpen(f) }

// Set the options for Driver.OpenConnector.
func DriverOpenConnector(f func(*StepOptions)) Option { return opts.DriverOpenConnector(f) }
