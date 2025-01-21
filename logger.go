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

func (x *logger) StepWithoutContext(step *StepOptions, fn func() error) error {
	return x.Step(context.Background(), step, fn)
}

func (x *logger) Step(ctx context.Context, step *StepOptions, fn func() error) error {
	x.Log(ctx, slog.Level(step.Start.Level), step.Start.Msg)
	t0 := time.Now()
	err := fn()
	lg := x.With(slog.Int64("duration", time.Since(t0).Nanoseconds()))
	var complete bool
	if step.ErrorHandler != nil {
		var attrs []slog.Attr
		complete, attrs = step.ErrorHandler(err)
		if len(attrs) > 0 {
			args := make([]interface{}, len(attrs))
			for i, attr := range attrs {
				args[i] = attr
			}
			lg = lg.With(args...)
		}
	} else {
		complete = err == nil
	}
	if !complete {
		lg.Log(ctx, slog.Level(step.Error.Level), step.Error.Msg, slog.Any("error", err))
	} else {
		lg.Log(ctx, slog.Level(step.Complete.Level), step.Complete.Msg)
	}
	return err
}
