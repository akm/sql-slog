package public

import (
	"errors"
	"testing"
)

func TestHandleRowsNextError(t *testing.T) {
	t.Parallel()
	complete, attrs := HandleRowsNextError(errors.New("dummy"))
	if complete {
		t.Fatal("Expected false")
	}
	if attrs != nil {
		t.Fatal("Expected nil")
	}
}
