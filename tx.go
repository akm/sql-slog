package sqlslog

import (
	"database/sql/driver"
)

func wrapTx(original driver.Tx, logger *logger) *txWrapper {
	return &txWrapper{original: original, logger: logger}
}

type txWrapper struct {
	original driver.Tx
	logger   *logger
}

var _ driver.Tx = (*txWrapper)(nil)

// Commit implements driver.Tx.
func (t *txWrapper) Commit() error {
	return t.logger.logAction("Tx.Commit", t.original.Commit)
}

// Rollback implements driver.Tx.
func (t *txWrapper) Rollback() error {
	return t.logger.logAction("Tx.Rollback", t.original.Rollback)
}
