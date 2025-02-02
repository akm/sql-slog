package opts

import (
	"errors"
	"io"
	"testing"
)

func TestDriverOpenErrorHandler(t *testing.T) {
	t.Parallel()
	t.Run("postgres", func(t *testing.T) {
		t.Parallel()
		t.Run("unexpected error", func(t *testing.T) {
			t.Parallel()
			eh := DriverOpenErrorHandler("postgres")
			completed, attrs := eh(errors.New("unexpected error"))
			if completed {
				t.Error("expected completed to be false")
			}
			if attrs != nil {
				t.Error("expected attrs to be nil")
			}
		})
		t.Run("io.EOF", func(t *testing.T) {
			t.Parallel()
			eh := DriverOpenErrorHandler("postgres")
			completed, attrs := eh(io.EOF)
			if !completed {
				t.Errorf("expected completed to be true")
			}
			if attrs == nil {
				t.Error("expected attrs to be non-nil")
			}
		})
	})
}
