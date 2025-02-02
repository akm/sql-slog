package opts

import (
	"log/slog"
	"testing"
	"time"
)

func TestDurationAttrFunc(t *testing.T) {
	t.Run("unexpected DurationType", func(t *testing.T) {
		f := DurationAttrFunc("key", DurationType(999))
		attr := f(123 * time.Nanosecond)
		if attr.Key != "key" {
			t.Errorf("expected: key, but got %s", attr.Key)
		}
		if attr.Value.Kind() != slog.KindInt64 {
			t.Errorf("expected: KindInt64, but got %s", attr.Value.Kind())
		}
		if attr.Value.Int64() != 123 {
			t.Errorf("expected: 123, but got %d", attr.Value.Int64())
		}
	})
}

func TestDurationString(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		gen      func() time.Duration
		expected string
	}{
		{
			gen:      func() time.Duration { return time.Duration(0) },
			expected: "0s",
		},
		{
			gen:      func() time.Duration { return time.Duration(1) },
			expected: "1ns",
		},
		{
			gen:      func() time.Duration { return time.Duration(1e3) },
			expected: "1µs",
		},
		{
			gen:      func() time.Duration { return time.Duration(1e6) },
			expected: "1ms",
		},
		{
			gen:      func() time.Duration { return time.Duration(1e9) },
			expected: "1s",
		},
		{
			gen:      func() time.Duration { return time.Duration(1e9 + 1) },
			expected: "1.000000001s",
		},
		{
			gen:      func() time.Duration { return time.Duration(1e9 + 1e3) },
			expected: "1.000001s",
		},
		{
			gen:      func() time.Duration { return time.Duration(1e9 * 60) },
			expected: "1m0s",
		},
		{
			gen:      func() time.Duration { return time.Duration(1e9*60 + 1) },
			expected: "1m0.000000001s",
		},
		{
			gen:      func() time.Duration { return time.Duration(1e9 * 60 * 60) },
			expected: "1h0m0s",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.expected, func(t *testing.T) {
			t.Parallel()
			d := tc.gen()
			actual := d.String()
			if actual != tc.expected {
				t.Errorf("expected: %s, but got %s", tc.expected, actual)
			}
		})
	}
}
