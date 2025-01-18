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

func (x *logger) logAction(proc *ProcOptions, fn func() error) error {
	ctx := context.Background()
	x.Log(ctx, slog.Level(proc.Start.level), proc.Start.name)
	t0 := time.Now()
	err := fn()
	lg := x.With(slog.Int64("duration", time.Since(t0).Nanoseconds()))
	if err != nil {
		lg.Log(ctx, slog.Level(proc.Error.level), proc.Error.name, "error", err)
		return err
	}
	lg.Log(ctx, slog.Level(proc.Complete.level), proc.Complete.name)
	return nil
}

func (x *logger) logActionContext(ctx context.Context, proc *ProcOptions, fn func() error) error {
	x.Log(ctx, slog.Level(proc.Start.level), proc.Start.name)
	t0 := time.Now()
	err := fn()
	lg := x.With(slog.Int64("duration", time.Since(t0).Nanoseconds()))
	if err != nil {
		lg.Log(ctx, slog.Level(proc.Error.level), proc.Error.name, "error", err)
		return err
	}
	lg.Log(ctx, slog.Level(proc.Complete.level), proc.Complete.name)
	return nil
}
