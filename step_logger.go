package sqlslog

import (
	"log/slog"
	"time"

	"github.com/akm/sql-slog/internal/wrap"
)

type StepLogger = wrap.StepLogger

func NewStepLogger(rawLogger *slog.Logger, durationAttr func(d time.Duration) slog.Attr) *StepLogger {
	return wrap.NewStepLogger(rawLogger, durationAttr)
}

func WithNilAttr(f func() error) func() (*slog.Attr, error) {
	return wrap.WithNilAttr(f)
}

func IgnoreAttr(att *slog.Attr, err error) error {
	return wrap.IgnoreAttr(att, err)
}
