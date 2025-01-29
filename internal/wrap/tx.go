package wrap

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
	return IgnoreAttr(t.logger.StepWithoutContext(&t.logger.Options.TxCommit, WithNilAttr(t.original.Commit)))
}

// Rollback implements driver.Tx.
func (t *txWrapper) Rollback() error {
	return IgnoreAttr(t.logger.StepWithoutContext(&t.logger.Options.TxRollback, WithNilAttr(t.original.Rollback)))
}
