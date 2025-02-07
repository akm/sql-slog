package sqlslog

import "testing"

func TestStepOptionsSetLevel(t *testing.T) {
	t.Parallel()

	newOpt := func(start, err, comp Level) *StepOptions {
		return &StepOptions{
			Start:    EventOptions{Level: start},
			Error:    EventOptions{Level: err},
			Complete: EventOptions{Level: comp},
		}
	}

	tests := []struct {
		name string
		o    *StepOptions
		lv   Level
		want *StepOptions
	}{
		{
			name: "LevelTrace",
			o:    newOpt(LevelDebug, LevelError, LevelInfo),
			lv:   LevelTrace,
			want: newOpt(LevelVerbose, LevelError, LevelTrace),
		},
		{
			name: "LevelDebug",
			o:    newOpt(LevelInfo, LevelError, LevelInfo),
			lv:   LevelDebug,
			want: newOpt(LevelTrace, LevelError, LevelDebug),
		},
		{
			name: "LevelInfo",
			o:    newOpt(LevelTrace, LevelError, LevelDebug),
			lv:   LevelInfo,
			want: newOpt(LevelDebug, LevelError, LevelInfo),
		},
		{
			name: "LevelWarn",
			o:    newOpt(LevelTrace, LevelError, LevelDebug),
			lv:   LevelWarn,
			want: newOpt(LevelInfo, LevelError, LevelWarn),
		},
		{
			name: "LevelError",
			o:    newOpt(LevelTrace, LevelError, LevelDebug),
			lv:   LevelError,
			want: newOpt(LevelWarn, LevelError, LevelError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.o.SetLevel(tt.lv)
			if !tt.o.compare(tt.want) {
				t.Errorf("StepOptions.SetLevel() = %v, want %v", tt.o, tt.want)
			}
			if tt.o.Complete.Level != tt.lv {
				t.Errorf("StepOptions.SetLevel() = %v, want %v", tt.o.Complete.Level, tt.lv)
			}
		})
	}
}

func TestDefaultStepOptions(t *testing.T) {
	t.Parallel()
	t.Run("LevelError", func(t *testing.T) {
		t.Parallel()
		o := defaultStepOptions(StepEventMsgWithoutEventName, Step("test"), LevelError)
		if o.Start.Level != LevelInfo {
			t.Errorf("Expected %v, but got %v", LevelInfo, o.Start.Level)
		}
		if o.Complete.Level != LevelError {
			t.Errorf("Expected %v, but got %v", LevelError, o.Complete.Level)
		}
	})
}
