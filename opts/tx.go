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
