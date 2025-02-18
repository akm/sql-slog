package sqlslog

import (
	"fmt"
	"testing"
)

func TestLevelString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		baseName string
		base     Level
		diff     Level
		want     string
	}{
		{baseName: "LevelVerbose", base: LevelVerbose, diff: -1, want: "VERBOSE-1"},
		{baseName: "LevelVerbose", base: LevelVerbose, diff: 0, want: "VERBOSE"},
		{baseName: "LevelVerbose", base: LevelVerbose, diff: +3, want: "VERBOSE+3"},
		{baseName: "LevelTrace", base: LevelTrace, diff: -1, want: "VERBOSE+3"},
		{baseName: "LevelTrace", base: LevelTrace, diff: 0, want: "TRACE"},
		{baseName: "LevelTrace", base: LevelTrace, diff: 3, want: "TRACE+3"},
		{baseName: "LevelDebug", base: LevelDebug, diff: -1, want: "TRACE+3"},
		{baseName: "LevelDebug", base: LevelDebug, diff: 0, want: "DEBUG"},
		{baseName: "LevelDebug", base: LevelDebug, diff: +3, want: "DEBUG+3"},
		{baseName: "LevelInfo", base: LevelInfo, diff: -1, want: "DEBUG+3"},
		{baseName: "LevelInfo", base: LevelInfo, diff: 0, want: "INFO"},
		{baseName: "LevelInfo", base: LevelInfo, diff: +1, want: "INFO+1"},
		{baseName: "LevelInfo", base: LevelInfo, diff: +3, want: "INFO+3"},
		{baseName: "LevelWarn", base: LevelWarn, diff: -1, want: "INFO+3"},
		{baseName: "LevelWarn", base: LevelWarn, diff: 0, want: "WARN"},
		{baseName: "LevelWarn", base: LevelWarn, diff: +3, want: "WARN+3"},
		{baseName: "LevelError", base: LevelError, diff: -1, want: "WARN+3"},
		{baseName: "LevelError", base: LevelError, diff: 0, want: "ERROR"},
		{baseName: "LevelError", base: LevelError, diff: +1, want: "ERROR+1"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s%d", tt.baseName, tt.diff), func(t *testing.T) {
			t.Parallel()
			if got := (tt.base + tt.diff).String(); got != tt.want {
				t.Errorf("Level.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	t.Parallel()
	t.Run("Valid case", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name string
			in   string
			want Level
		}{
			{name: "LevelVerbose", in: "VERBOSE", want: LevelVerbose},
			{name: "LevelTrace", in: "TRACE", want: LevelTrace},
			{name: "LevelDebug", in: "DEBUG", want: LevelDebug},
			{name: "LevelInfo", in: "INFO", want: LevelInfo},
			{name: "LevelWarn", in: "WARN", want: LevelWarn},
			{name: "LevelError", in: "ERROR", want: LevelError},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				if got, err := ParseLevel(tt.in); err != nil || got != tt.want {
					t.Errorf("ParseLevel() = %v, %v, want %v, true", got, err, tt.want)
				}
			})
		}
	})
	t.Run("Unknown level", func(t *testing.T) {
		t.Parallel()
		tests := []string{"", "UNKNOWN", "TRACE+", "TRACE-1", "TRACE+1", "TRACE+1+1"}
		for _, in := range tests {
			t.Run(in, func(t *testing.T) {
				t.Parallel()
				lv, err := ParseLevel(in)
				if err == nil {
					t.Fatal("Expected non-nil")
				}
				if lv != 0 {
					t.Fatalf("Expected 0, got %v", lv)
				}
			})
		}
	})
}

func TestParseLevelWithDefault(t *testing.T) {
	t.Parallel()
	t.Run("Valid case", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name string
			in   string
			want Level
		}{
			{name: "LevelVerbose", in: "VERBOSE", want: LevelVerbose},
			{name: "LevelTrace", in: "TRACE", want: LevelTrace},
			{name: "LevelDebug", in: "DEBUG", want: LevelDebug},
			{name: "LevelInfo", in: "INFO", want: LevelInfo},
			{name: "LevelWarn", in: "WARN", want: LevelWarn},
			{name: "LevelError", in: "ERROR", want: LevelError},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				if got := ParseLevelWithDefault(tt.in, LevelError); got != tt.want {
					t.Errorf("ParseLevelWithDefault() = %v, want %v", got, tt.want)
				}
			})
		}
	})
	t.Run("Unknown level", func(t *testing.T) {
		t.Parallel()
		tests := []string{"", "UNKNOWN", "TRACE+", "TRACE-1", "TRACE+1", "TRACE+1+1"}
		for _, in := range tests {
			t.Run(in, func(t *testing.T) {
				t.Parallel()
				lv := ParseLevelWithDefault(in, LevelError)
				if lv != LevelError {
					t.Fatalf("Expected LevelError, got %v", lv)
				}
			})
		}
	})
}
