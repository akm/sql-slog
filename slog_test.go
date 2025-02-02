package sqlslog

import (
	"bytes"
	"testing"
)

func TestNewTextHandler(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		buf := bytes.NewBuffer(nil)
		h := NewTextHandler(buf, nil)
		if h == nil {
			t.Error("NewTextHandler() should return a handler")
		}
	})
}
