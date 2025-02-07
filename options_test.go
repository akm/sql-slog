package sqlslog

import "testing"

func TestSetStepEventMsgBuilder(t *testing.T) {
	t.Parallel()
	t.Run("defaultOptions", func(t *testing.T) {
		t.Parallel()
		opts := newOptions("dummy")
		if opts.DriverOptions.ConnOptions.Begin.Complete.Msg != "Conn.Begin" {
			t.Errorf("unexpected default value: %s", opts.DriverOptions.ConnOptions.Begin.Complete.Msg)
		}
	})

	t.Run("CustomStepEventMsgBuilder", func(t *testing.T) {
		t.Parallel()
		builder, backup := StepEventMsgWithEventName, stepEventMsgBuilder
		SetStepEventMsgBuilder(builder)
		defer SetStepEventMsgBuilder(backup)
		opts := newOptions("dummy")
		if opts.DriverOptions.ConnOptions.Begin.Complete.Msg != "Conn.Begin Complete" {
			t.Errorf("unexpected default value: %s", opts.DriverOptions.ConnOptions.Begin.Complete.Msg)
		}
	})
}
