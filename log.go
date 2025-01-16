package sqlslog

import (
	"context"
	"log/slog"
	"time"
)

func logAction(logger *slog.Logger, action string, fn func() error) error {
	logger.Debug(action + " Start")
	t0 := time.Now()
	err := fn()
	lg := logger.With(slog.Duration("duration", time.Since(t0)))
	if err != nil {
		lg.Error(action+" Error", "error", err)
		return err
	}
	lg.Info(action + " Complete")
	return nil
}

func logActionContext(ctx context.Context, logger *slog.Logger, action string, fn func() error) error {
	logger.DebugContext(ctx, action+" Start")
	t0 := time.Now()
	err := fn()
	lg := logger.With(slog.Duration("duration", time.Since(t0)))
	if err != nil {
		lg.ErrorContext(ctx, action+" Error", "error", err)
		return err
	}
	lg.InfoContext(ctx, action+" Complete")
	return nil
}
