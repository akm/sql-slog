package wrap

import (
	"context"
	"log/slog"
	"time"

	"github.com/akm/sql-slog/opts"
)

type logger struct {
	*slog.Logger
	Options *opts.Options
}

func newLogger(rawLogger *slog.Logger, opts *opts.Options) *logger {
	return &logger{
		Logger:  rawLogger,
		Options: opts,
	}
}

func (x *logger) With(kv ...interface{}) *logger {
	return newLogger(x.Logger.With(kv...), x.Options)
}

func (x *logger) StepWithoutContext(step *opts.StepOptions, fn func() (*slog.Attr, error)) (*slog.Attr, error) {
	return x.Step(context.Background(), step, fn)
}

func (x *logger) Step(ctx context.Context, step *opts.StepOptions, fn func() (*slog.Attr, error)) (*slog.Attr, error) {
	x.Log(ctx, slog.Level(step.Start.Level), step.Start.Msg)
	t0 := time.Now()
	attr, err := fn()
	lg := x.With(x.durationAttr(time.Since(t0)))
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
	switch {
	case !complete:
		lg.Log(ctx, slog.Level(step.Error.Level), step.Error.Msg, slog.Any("error", err))
	case attr != nil:
		lg.Log(ctx, slog.Level(step.Complete.Level), step.Complete.Msg, *attr)
	default:
		lg.Log(ctx, slog.Level(step.Complete.Level), step.Complete.Msg)
	}
	return attr, err
}

func (x *logger) durationAttr(d time.Duration) slog.Attr {
	key := x.Options.DurationKey
	switch x.Options.DurationType {
	case opts.DurationNanoSeconds:
		return slog.Int64(key, d.Nanoseconds())
	case opts.DurationMicroSeconds:
		return slog.Int64(key, d.Microseconds())
	case opts.DurationMilliSeconds:
		return slog.Int64(key, d.Milliseconds())
	case opts.DurationGoDuration:
		return slog.Duration(key, d)
	case opts.DurationString:
		return slog.String(key, d.String())
	default:
		return slog.Int64(key, d.Nanoseconds())
	}
}

func WithNilAttr(f func() error) func() (*slog.Attr, error) {
	return func() (*slog.Attr, error) {
		return nil, f()
	}
}

func IgnoreAttr(_ *slog.Attr, err error) error {
	return err
}
