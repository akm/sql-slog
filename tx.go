package sqlslog

import (
	"database/sql/driver"
	"log/slog"
)

func wrapTx(original driver.Tx, logger *slog.Logger) *txWrapper {
	return &txWrapper{original: original, logger: logger}
}

type txWrapper struct {
	original driver.Tx
	logger   *slog.Logger
}

var _ driver.Tx = (*txWrapper)(nil)

// Commit implements driver.Tx.
func (t *txWrapper) Commit() error {
	lg := t.logger
	lg.Debug("Commit Start")
	if err := t.original.Commit(); err != nil {
		lg.Error("Commit Error", "error", err)
		return err
	}
	lg.Info("Commit Complete")
	return nil
}

// Rollback implements driver.Tx.
func (t *txWrapper) Rollback() error {
	lg := t.logger
	lg.Debug("Rollback Start")
	if err := t.original.Rollback(); err != nil {
		lg.Error("Rollback Error", "error", err)
		return err
	}
	lg.Info("Rollback Complete")
	return nil
}
