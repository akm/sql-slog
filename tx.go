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
	return logAction(t.logger, "Tx.Commit", t.original.Commit)
}

// Rollback implements driver.Tx.
func (t *txWrapper) Rollback() error {
	return logAction(t.logger, "Tx.Rollback", t.original.Rollback)
}
