package sqlslog

import (
	"context"
	"log/slog"
	"time"
)

type logger struct {
	*slog.Logger
	options *options
}

func newLogger(rawLogger *slog.Logger, opts *options) *logger {
	return &logger{
		Logger:  rawLogger,
		options: opts,
	}
}

func (x *logger) With(kv ...interface{}) *logger {
	return newLogger(x.Logger.With(kv...), x.options)
}

func (x *logger) logAction(action string, fn func() error) error {
	x.Debug(action + " Start")
	t0 := time.Now()
	err := fn()
	lg := x.With(slog.Int64("duration", time.Since(t0).Nanoseconds()))
	if err != nil {
		lg.Error(action+" Error", "error", err)
		return err
	}
	lg.Info(action + " Complete")
	return nil
}

func (x *logger) logActionContext(ctx context.Context, action string, fn func() error) error {
	x.DebugContext(ctx, action+" Start")
	t0 := time.Now()
	err := fn()
	lg := x.With(slog.Int64("duration", time.Since(t0).Nanoseconds()))
	if err != nil {
		lg.ErrorContext(ctx, action+" Error", "error", err)
		return err
	}
	lg.InfoContext(ctx, action+" Complete")
	return nil
}
