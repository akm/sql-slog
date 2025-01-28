package sqlslog

// Duration is an option to specify duration value in log.
// The default is DurationNanoSeconds.
func Duration(v DurationType) Option {
	return func(o *Options) {
		o.DurationType = v
	}
}

// DurationKey is an option to specify the key for duration value in log.
// The default is specified by DurationKeyDefault.
func DurationKey(key string) Option {
	return func(o *Options) {
		o.DurationKey = key
	}
}
