package sqlslog

import "log/slog"

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
