package sqlslog

import (
	"database/sql/driver"
)

type txOptions struct {
	Commit   StepOptions
	Rollback StepOptions
}

func defaultTxOptions(msgb StepEventMsgBuilder) *txOptions {
	return &txOptions{
		Commit:   *defaultStepOptions(msgb, StepTxCommit, LevelInfo),
		Rollback: *defaultStepOptions(msgb, StepTxRollback, LevelInfo),
	}
}

func wrapTx(original driver.Tx, logger *stepLogger, options *txOptions) *txWrapper {
	return &txWrapper{original: original, logger: logger, options: options}
}

type txWrapper struct {
	original driver.Tx
	logger   *stepLogger
	options  *txOptions
}

var _ driver.Tx = (*txWrapper)(nil)

// Commit implements driver.Tx.
func (t *txWrapper) Commit() error {
	return ignoreAttr(t.logger.StepWithoutContext(&t.options.Commit, withNilAttr(t.original.Commit)))
}

// Rollback implements driver.Tx.
func (t *txWrapper) Rollback() error {
	return ignoreAttr(t.logger.StepWithoutContext(&t.options.Rollback, withNilAttr(t.original.Rollback)))
}
