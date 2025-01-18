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

func (x *logger) logAction(proc *StepOptions, fn func() error) error {
	ctx := context.Background()
	x.Log(ctx, slog.Level(proc.Start.Level), proc.Start.Msg)
	t0 := time.Now()
	err := fn()
	lg := x.With(slog.Int64("duration", time.Since(t0).Nanoseconds()))
	if err != nil {
		lg.Log(ctx, slog.Level(proc.Error.Level), proc.Error.Msg, "error", err)
		return err
	}
	lg.Log(ctx, slog.Level(proc.Complete.Level), proc.Complete.Msg)
	return nil
}

func (x *logger) logActionContext(ctx context.Context, proc *StepOptions, fn func() error) error {
	x.Log(ctx, slog.Level(proc.Start.Level), proc.Start.Msg)
	t0 := time.Now()
	err := fn()
	lg := x.With(slog.Int64("duration", time.Since(t0).Nanoseconds()))
	if err != nil {
		lg.Log(ctx, slog.Level(proc.Error.Level), proc.Error.Msg, "error", err)
		return err
	}
	lg.Log(ctx, slog.Level(proc.Complete.Level), proc.Complete.Msg)
	return nil
}
