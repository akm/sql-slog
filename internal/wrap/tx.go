package wrap

import (
	"database/sql/driver"

	"github.com/akm/sql-slog/opts"
)

type TxOptions = opts.TxOptions

func WrapTx(original driver.Tx, logger *StepLogger, options *TxOptions) driver.Tx { // nolint:revive
	return &txWrapper{original: original, logger: logger, options: options}
}

type txWrapper struct {
	original driver.Tx
	logger   *StepLogger
	options  *TxOptions
}

var _ driver.Tx = (*txWrapper)(nil)

// Commit implements driver.Tx.
func (t *txWrapper) Commit() error {
	return IgnoreAttr(t.logger.StepWithoutContext(t.options.Commit, WithNilAttr(t.original.Commit)))
}

// Rollback implements driver.Tx.
func (t *txWrapper) Rollback() error {
	return IgnoreAttr(t.logger.StepWithoutContext(t.options.Rollback, WithNilAttr(t.original.Rollback)))
}
