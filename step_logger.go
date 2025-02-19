package sqlslog

import (
	"context"
	"log/slog"
	"time"
)

type stepLoggerOptions struct {
	durationKey  string
	durationType DurationType
}

func defaultStepLoggerOptions() stepLoggerOptions {
	return stepLoggerOptions{
		durationKey:  DurationKeyDefault,
		durationType: DurationNanoSeconds,
	}
}

type stepLogger struct {
	*slog.Logger
	durationAttr func(d time.Duration) slog.Attr
}

func newStepLogger(logger *slog.Logger, opts stepLoggerOptions) *stepLogger {
	return &stepLogger{
		Logger:       logger,
		durationAttr: durationAttrFunc(opts.durationKey, opts.durationType),
	}
}

func (x *stepLogger) With(kv ...interface{}) *stepLogger {
	return &stepLogger{
		Logger:       x.Logger.With(kv...),
		durationAttr: x.durationAttr,
	}
}

func (x *stepLogger) StepWithoutContext(step *StepOptions, fn func() (*slog.Attr, error)) (*slog.Attr, error) {
	return x.Step(context.Background(), step, fn)
}

func (x *stepLogger) Step(ctx context.Context, step *StepOptions, fn func() (*slog.Attr, error)) (*slog.Attr, error) {
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

func durationAttrFunc(key string, dt DurationType) func(d time.Duration) slog.Attr {
	switch dt {
	case DurationNanoSeconds:
		return func(d time.Duration) slog.Attr { return slog.Int64(key, d.Nanoseconds()) }
	case DurationMicroSeconds:
		return func(d time.Duration) slog.Attr { return slog.Int64(key, d.Microseconds()) }
	case DurationMilliSeconds:
		return func(d time.Duration) slog.Attr { return slog.Int64(key, d.Milliseconds()) }
	case DurationGoDuration:
		return func(d time.Duration) slog.Attr { return slog.Duration(key, d) }
	case DurationString:
		return func(d time.Duration) slog.Attr { return slog.String(key, d.String()) }
	default:
		return func(d time.Duration) slog.Attr { return slog.Int64(key, d.Nanoseconds()) }
	}
}

func withNilAttr(f func() error) func() (*slog.Attr, error) {
	return func() (*slog.Attr, error) {
		return nil, f()
	}
}

func ignoreAttr(_ *slog.Attr, err error) error {
	return err
}
