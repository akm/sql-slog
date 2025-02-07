package sqlslog

import "testing"

func TestEventString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		e    Event
		want string
	}{
		{
			name: "EventStart",
			e:    EventStart,
			want: "Start",
		},
		{
			name: "EventError",
			e:    EventError,
			want: "Error",
		},
		{
			name: "EventComplete",
			e:    EventComplete,
			want: "Complete",
		},
		{
			name: "Unknown",
			e:    Event(0),
			want: "Unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Event.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
