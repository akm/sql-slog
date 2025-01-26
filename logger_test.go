package sqlslog

import (
	"log/slog"
	"testing"
	"time"
)

func TestLoggerDurationAttr(t *testing.T) {
	t.Parallel()
	key := "duration"
	testcases := []struct {
		value        time.Duration
		durationType DurationType
		expected     func(t *testing.T, attr slog.Attr)
	}{
		{
			value:        time.Duration(1),
			durationType: DurationNanoSeconds,
			expected: func(t *testing.T, attr slog.Attr) {
				t.Helper()
				if v, ok := attr.Value.Any().(int64); !ok {
					t.Errorf("expected: %T, but got %T", int64(0), v)
				} else if v != 1 {
					t.Errorf("expected: %d, but got %d", 1, v)
				}
			},
		},
		{
			value:        time.Duration(2_000),
			durationType: DurationMicroSeconds,
			expected: func(t *testing.T, attr slog.Attr) {
				t.Helper()
				if v, ok := attr.Value.Any().(int64); !ok {
					t.Errorf("expected: %T, but got %T", int64(0), v)
				} else if v != 2 {
					t.Errorf("expected: %d, but got %d", 2, v)
				}
			},
		},
		{
			value:        time.Duration(3_000_000),
			durationType: DurationMilliSeconds,
			expected: func(t *testing.T, attr slog.Attr) {
				t.Helper()
				if v, ok := attr.Value.Any().(int64); !ok {
					t.Errorf("expected: %T, but got %T", int64(0), v)
				} else if v != 3 {
					t.Errorf("expected: %d, but got %d", 3, v)
				}
			},
		},
		{
			value:        time.Duration(4_000_000_000),
			durationType: DurationGoDuration,
			expected: func(t *testing.T, attr slog.Attr) {
				t.Helper()
				if attr.Value.Duration() != time.Duration(4_000_000_000) {
					t.Errorf("expected: %d, but got %d", 4_000_000_000, attr.Value.Duration())
				}
			},
		},
		{
			value:        time.Duration(567_000_000),
			durationType: DurationString,
			expected: func(t *testing.T, attr slog.Attr) {
				t.Helper()
				if v, ok := attr.Value.Any().(string); !ok {
					t.Errorf("expected: %T, but got %T", "", v)
				} else if v != "567ms" {
					t.Errorf("expected: %s, but got %s", "567ms", v)
				}
			},
		},
		{
			value:        time.Duration(890_000),
			durationType: DurationType(-1),
			expected: func(t *testing.T, attr slog.Attr) {
				t.Helper()
				if v, ok := attr.Value.Any().(int64); !ok {
					t.Errorf("expected: %T, but got %T", int64(0), v)
				} else if v != 890_000 {
					t.Errorf("expected: %d, but got %d", 890_000, v)
				}
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.value.String(), func(t *testing.T) {
			t.Parallel()
			attr := newLogger(nil, &options{durationKey: key, durationType: tc.durationType}).durationAttr(tc.value)
			tc.expected(t, attr)
		})
	}
}
