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
