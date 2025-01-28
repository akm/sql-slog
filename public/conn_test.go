package public

import (
	"errors"
	"testing"
)

func TestConnExecContextErrorHandler(t *testing.T) {
	t.Parallel()
	errHandler := ConnExecContextErrorHandler("mysql")
	complete, attrs := errHandler(errors.New("dummy"))
	if complete {
		t.Fatal("Expected false")
	}
	if attrs != nil {
		t.Fatal("Expected nil")
	}
}

func TestConnQueryContextErrorHandler(t *testing.T) {
	t.Parallel()
	t.Run("mysql", func(t *testing.T) {
		t.Parallel()
		errHandler := ConnQueryContextErrorHandler("mysql")
		t.Run("nil error", func(t *testing.T) {
			t.Parallel()
			complete, attrs := errHandler(nil)
			if !complete {
				t.Fatal("Expected true")
			}
			if attrs != nil {
				t.Fatal("Expected nil")
			}
		})
		t.Run("unexpected error", func(t *testing.T) {
			t.Parallel()
			complete, attrs := errHandler(errors.New("dummy"))
			if complete {
				t.Fatal("Expected false")
			}
			if attrs != nil {
				t.Fatal("Expected nil")
			}
		})
	})
}
