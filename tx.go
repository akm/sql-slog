package sqlslog

import (
	"database/sql/driver"
)

func wrapTx(original driver.Tx, logger *SqlLogger) *txWrapper {
	return &txWrapper{original: original, logger: logger}
}

type txWrapper struct {
	original driver.Tx
	logger   *SqlLogger
}

var _ driver.Tx = (*txWrapper)(nil)

// Commit implements driver.Tx.
func (t *txWrapper) Commit() error {
	return IgnoreAttr(t.logger.StepWithoutContext(&t.logger.options.TxCommit, WithNilAttr(t.original.Commit)))
}

// Rollback implements driver.Tx.
func (t *txWrapper) Rollback() error {
	return IgnoreAttr(t.logger.StepWithoutContext(&t.logger.options.TxRollback, WithNilAttr(t.original.Rollback)))
}
