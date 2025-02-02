package sqlslog

import (
	"testing"

	"github.com/akm/sql-slog/internal/opts"
)

func TestSetStepLogMsgFormatter(t *testing.T) {
	t.Parallel()
	t.Run("StepLogMsgWithEventName", func(t *testing.T) {
		t.Parallel()
		SetStepLogMsgFormatter(opts.StepLogMsgWithEventName)
		defer SetStepLogMsgFormatter(opts.StepLogMsgWithoutEventName)

		options := opts.NewOptions("dummy")
		if options.OpenOptions.Open.Complete.Msg != "Open Complete" {
			t.Errorf("unexpected default value: %s", options.OpenOptions.Open.Complete.Msg)
		}
	})
}
