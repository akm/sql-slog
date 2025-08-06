package sqlslog

import (
	"database/sql/driver"
	"testing"
	"time"
)

func TestFormatNamedValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []driver.NamedValue
		want string
	}{
		{
			name: "empty",
			args: nil,
			want: "[]",
		},
		{
			name: "no names without ordinals",
			args: []driver.NamedValue{{Value: "A"}, {Value: 2}},
			want: "[\"A\",2]",
		},
		{
			name: "no names with sorted ordinals",
			args: []driver.NamedValue{{Ordinal: 1, Value: "a"}, {Ordinal: 2, Value: 100}, {Ordinal: 3, Value: 1.234}},
			want: "[\"a\",100,1.234]",
		},
		{
			name: "no names with unsorted ordinals",
			args: []driver.NamedValue{{Ordinal: 2, Value: 100}, {Ordinal: 3, Value: 1.234}, {Ordinal: 1, Value: "a"}},
			want: "[\"a\",100,1.234]",
		},
		{
			name: "no names with invalid ordinals",
			args: []driver.NamedValue{{Ordinal: 2, Value: 100}, {Ordinal: 5, Value: 1.234}, {Ordinal: 0, Value: "a"}},
			want: "[\"a\",100,1.234]",
		},
		{
			name: "all names",
			args: []driver.NamedValue{{Name: "a", Value: 1}, {Name: "b", Value: 2}},
			want: "{a:1,b:2}",
		},
		{
			name: "mixed names",
			args: []driver.NamedValue{{Ordinal: 0, Name: "a", Value: 1}, {Ordinal: 1, Name: "", Value: 2}},
			want: "[[0]a:1,[1]:2]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := formatNamedValues(tt.args)
			if got != tt.want {
				t.Errorf("formatNamedValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value driver.Value
		want  string
	}{
		{value: nil, want: "<nil>"},
		{value: true, want: "true"},
		{value: false, want: "false"},
		{value: 123, want: "123"},
		{value: 1.23, want: "1.23"},
		{value: "test", want: "\"test\""},
		{value: []byte("bytes"), want: "\"bytes\""},
		{value: []int{1, 2, 3}, want: "[1 2 3]"},
		{value: map[string]int{"a": 1, "b": 2}, want: "map[a:1 b:2]"},
		{value: struct{ A int }{A: 1}, want: "{1}"},
		{value: time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC), want: "2023-10-01 12:00:00 +0000 UTC"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			got := formatValue(tt.value)
			if got != tt.want {
				t.Errorf("formatValue(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
