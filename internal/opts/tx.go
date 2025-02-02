package opts

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

// Set the options for Tx.Commit.
func TxCommit(f func(*StepOptions)) Option { return func(o *Options) { f(&o.TxCommit) } }

// Set the options for Tx.Rollback.
func TxRollback(f func(*StepOptions)) Option { return func(o *Options) { f(&o.TxRollback) } }
