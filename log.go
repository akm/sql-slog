package sqlslog

import (
	"context"
	"log/slog"
)

func logAction(logger *slog.Logger, action string, fn func() error) error {
	logger.Debug(action + " Start")
	err := fn()
	if err != nil {
		logger.Error(action+" Error", "error", err)
		return err
	}
	logger.Info(action + " Complete")
	return nil
}

func logActionContext(ctx context.Context, logger *slog.Logger, action string, fn func() error) error {
	logger.DebugContext(ctx, action+" Start")
	err := fn()
	if err != nil {
		logger.ErrorContext(ctx, action+" Error", "error", err)
		return err
	}
	logger.InfoContext(ctx, action+" Complete")
	return nil
}
