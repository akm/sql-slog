package sqlslog

import (
	"database/sql/driver"
	"testing"
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
