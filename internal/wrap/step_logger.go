package wrap

import (
	"context"
	"log/slog"
	"time"
)

type StepLogger struct {
	*slog.Logger
	durationAttr func(d time.Duration) slog.Attr
}

func NewStepLogger(rawLogger *slog.Logger, durationAttr func(d time.Duration) slog.Attr) *StepLogger {
	return &StepLogger{
		Logger:       rawLogger,
		durationAttr: durationAttr,
	}
}

func (x *StepLogger) With(kv ...interface{}) *StepLogger {
	return NewStepLogger(x.Logger.With(kv...), x.durationAttr)
}

func (x *StepLogger) StepWithoutContext(step *StepOptions, fn func() (*slog.Attr, error)) (*slog.Attr, error) {
	return x.Step(context.Background(), step, fn)
}

func (x *StepLogger) Step(ctx context.Context, step *StepOptions, fn func() (*slog.Attr, error)) (*slog.Attr, error) {
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

func WithNilAttr(f func() error) func() (*slog.Attr, error) {
	return func() (*slog.Attr, error) {
		return nil, f()
	}
}

func IgnoreAttr(_ *slog.Attr, err error) error {
	return err
}
