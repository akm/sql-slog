package sqlslog

import (
	"database/sql/driver"
)

func wrapTx(original driver.Tx, logger *stepLogger) *txWrapper {
	return &txWrapper{original: original, logger: logger}
}

type txWrapper struct {
	original driver.Tx
	logger   *stepLogger
}

var _ driver.Tx = (*txWrapper)(nil)

// Commit implements driver.Tx.
func (t *txWrapper) Commit() error {
	return ignoreAttr(t.logger.StepWithoutContext(&t.logger.options.txCommit, withNilAttr(t.original.Commit)))
}

// Rollback implements driver.Tx.
func (t *txWrapper) Rollback() error {
	return ignoreAttr(t.logger.StepWithoutContext(&t.logger.options.txRollback, withNilAttr(t.original.Rollback)))
}
