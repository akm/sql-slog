package sqlslog

import (
	"database/sql/driver"
	"fmt"
	"slices"
	"strings"
)

func cmpNamedValueByOrdinal(a, b driver.NamedValue) int {
	return a.Ordinal - b.Ordinal
}

func formatNamedValues(args []driver.NamedValue) string {
	if len(args) == 0 {
		return "[]"
	}
	if !slices.ContainsFunc(args, func(arg driver.NamedValue) bool { return arg.Name == "" }) {
		var b strings.Builder
		b.WriteString("{")
		for i, arg := range args {
			if i > 0 {
				b.WriteString(",")
			}
			fmt.Fprintf(&b, "%s:", arg.Name)
			b.WriteString(formatValue(arg.Value))
		}
		b.WriteString("}")
		return b.String()
	}
	if !slices.ContainsFunc(args, func(arg driver.NamedValue) bool { return arg.Name != "" }) {
		if !slices.IsSortedFunc(args, cmpNamedValueByOrdinal) {
			slices.SortFunc(args, cmpNamedValueByOrdinal)
		}
		var b strings.Builder
		b.WriteString("[")
		for i, arg := range args {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(formatValue(arg.Value))
		}
		b.WriteString("]")
		return b.String()
	}
	var b strings.Builder
	b.WriteString("[")
	for i, arg := range args {
		if i > 0 {
			b.WriteString(",")
		}
		fmt.Fprintf(&b, "[%d]%s:", arg.Ordinal, arg.Name)
		b.WriteString(formatValue(arg.Value))
	}
	b.WriteString("]")
	return b.String()
}

func formatValues(values []driver.Value) string {
	if len(values) == 0 {
		return "[]"
	}
	var b strings.Builder
	b.WriteString("[")
	for i, value := range values {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(formatValue(value))
	}
	b.WriteString("]")
	return b.String()
}

func formatValue(v any) string {
	switch v := v.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case []byte:
		return fmt.Sprintf("%q", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
