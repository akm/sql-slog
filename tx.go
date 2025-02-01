package sqlslog

import (
	"database/sql/driver"
)

type TxOptions struct {
	Commit   *StepOptions
	Rollback *StepOptions
}

func WrapTx(original driver.Tx, logger *logger, options *TxOptions) *txWrapper {
	return &txWrapper{original: original, logger: logger, options: options}
}

type txWrapper struct {
	original driver.Tx
	logger   *logger
	options  *TxOptions
}

var _ driver.Tx = (*txWrapper)(nil)

// Commit implements driver.Tx.
func (t *txWrapper) Commit() error {
	return ignoreAttr(t.logger.StepWithoutContext(t.options.Commit, withNilAttr(t.original.Commit)))
}

// Rollback implements driver.Tx.
func (t *txWrapper) Rollback() error {
	return ignoreAttr(t.logger.StepWithoutContext(t.options.Rollback, withNilAttr(t.original.Rollback)))
}
