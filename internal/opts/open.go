package opts

type OpenOptions struct {
	Open   *StepOptions
	Driver *DriverOptions
}

func DefaultOpenOptions(driverName string, formatter StepLogMsgFormatter) *OpenOptions {
	return &OpenOptions{
		Open:   DefaultStepOptions(formatter, "Open", LevelInfo),
		Driver: DefaultDriverOptions(driverName, formatter),
	}
}

// Set the options for sqlslog.Open.
func SqlslogOpen(f func(*StepOptions)) Option { return func(o *Options) { f(&o.SqlslogOpen) } } // nolint:revive
