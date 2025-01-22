package sqlslog

import (
	"log/slog"
	"testing"
	"time"
)

func TestLoggerDurationAttr(t *testing.T) {
	key := "duration"
	testcases := []struct {
		value        time.Duration
		durationType DurationType
		expection    func(t *testing.T, attr slog.Attr)
	}{
		{
			value:        time.Duration(1),
			durationType: DurationNanoSeconds,
			expection: func(t *testing.T, attr slog.Attr) {
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
			expection: func(t *testing.T, attr slog.Attr) {
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
			expection: func(t *testing.T, attr slog.Attr) {
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
			expection: func(t *testing.T, attr slog.Attr) {
				if attr.Value.Duration() != time.Duration(4_000_000_000) {
					t.Errorf("expected: %d, but got %d", 4_000_000_000, attr.Value.Duration())
				}
			},
		},
		{
			value:        time.Duration(567_000_000),
			durationType: DurationString,
			expection: func(t *testing.T, attr slog.Attr) {
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
			expection: func(t *testing.T, attr slog.Attr) {
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
			attr := newLogger(nil, &options{durationKey: key, durationType: tc.durationType}).durationAttr(tc.value)
			tc.expection(t, attr)
		})
	}
}
