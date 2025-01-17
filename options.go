package sqlslog

import "log/slog"

type options struct {
	logger *slog.Logger
}

type Option func(*options)

func newDefaultOptions() *options {
	return &options{
		logger: slog.Default(),
	}
}

func Logger(logger *slog.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func newOptions(opts ...Option) *options {
	o := newDefaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	return o
}
