package sqlslog

import "testing"

func TestSetStepEventMsgBuilder(t *testing.T) { // nolint:paralleltest
	t.Run("defaultOptions", func(t *testing.T) { // nolint:paralleltest
		opts := newOptions("dummy")
		if opts.DriverOptions.ConnOptions.Begin.Complete.Msg != "Conn.Begin" {
			t.Errorf("unexpected default value: %s", opts.DriverOptions.ConnOptions.Begin.Complete.Msg)
		}
	})

	t.Run("CustomStepEventMsgBuilder", func(t *testing.T) { // nolint:paralleltest
		builder, backup := StepEventMsgWithEventName, stepEventMsgBuilder
		SetStepEventMsgBuilder(builder)
		defer SetStepEventMsgBuilder(backup)
		opts := newOptions("dummy")
		if opts.DriverOptions.ConnOptions.Begin.Complete.Msg != "Conn.Begin Complete" {
			t.Errorf("unexpected default value: %s", opts.DriverOptions.ConnOptions.Begin.Complete.Msg)
		}
	})
}
