package sqlslog

import (
	"database/sql/driver"
)

type TxOptions struct {
	Commit   *StepOptions
	Rollback *StepOptions
}

func DefaultTxOptions(formatter StepLogMsgFormatter) *TxOptions {
	return &TxOptions{
		Commit:   DefaultStepOptions(formatter, "Tx.Commit", LevelInfo),
		Rollback: DefaultStepOptions(formatter, "Tx.Rollback", LevelInfo),
	}
}

func WrapTx(original driver.Tx, logger *logger, options *TxOptions) driver.Tx {
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
