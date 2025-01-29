package opts

import "testing"

func TestSetStepLogMsgFormatter(t *testing.T) {
	t.Parallel()
	t.Run("defaultOptions", func(t *testing.T) {
		t.Parallel()
		opts := NewOptions("dummy")
		if opts.ConnBegin.Complete.Msg != "Conn.Begin" {
			t.Errorf("unexpected default value: %s", opts.ConnBegin.Complete.Msg)
		}
	})

	t.Run("CustomStepLogMsgFormatter", func(t *testing.T) {
		t.Parallel()
		formatter, backup := StepLogMsgWithEventName, stepLogMsgFormatter
		SetStepLogMsgFormatter(formatter)
		defer SetStepLogMsgFormatter(backup)
		opts := NewOptions("dummy")
		if opts.ConnBegin.Complete.Msg != "Conn.Begin Complete" {
			t.Errorf("unexpected default value: %s", opts.ConnBegin.Complete.Msg)
		}
	})
}
