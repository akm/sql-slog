package testhelper

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type LogsAssertion struct {
	buf *bytes.Buffer
}

func NewLogAssertion(buf *bytes.Buffer) *LogsAssertion {
	return &LogsAssertion{buf: buf}
}

func (a *LogsAssertion) Start() {
	a.buf.Reset()
}

func (a *LogsAssertion) JsonLines(t *testing.T) []map[string]interface{} {
	return parseJsonLines(t, a.buf.Bytes())
}

func (a *LogsAssertion) Assert(t *testing.T, expected []map[string]interface{}) {
	actual := a.JsonLines(t)
	assertMapSlice(t, expected, actual,
		ignore("time"),
		when(msgEndsWith(" Complete"), deleteIfFloat64("duration")),
		when(msgEndsWith(" Error"), deleteIfFloat64("duration")),
	)
}

func (a *LogsAssertion) AssertEmpty(t *testing.T) {
	a.Assert(t, []map[string]interface{}{})
}

func parseJsonLines(t *testing.T, b []byte) []map[string]interface{} {
	t.Helper()
	lines := bytes.Split(b, []byte("\n"))
	results := []map[string]interface{}{}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		result := map[string]interface{}{}
		if err := json.Unmarshal(line, &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}
		results = append(results, result)
	}
	return results
}

type processor = func(t *testing.T, m map[string]interface{})

func ignore(fields ...string) processor {
	return func(_ *testing.T, m map[string]interface{}) {
		for _, f := range fields {
			delete(m, f)
		}
	}
}

func when(predicate func(map[string]interface{}) bool, processors ...processor) processor {
	return func(t *testing.T, m map[string]interface{}) {
		if predicate(m) {
			for _, p := range processors {
				p(t, m)
			}
		}
	}
}

func msgEndsWith(suffix string) func(map[string]interface{}) bool {
	return func(m map[string]interface{}) bool {
		msg, ok := m["msg"].(string)
		return ok && msg[len(msg)-len(suffix):] == suffix
	}
}

// Parsed value from JSON is float64
func deleteIfFloat64(key string) processor {
	return func(t *testing.T, m map[string]interface{}) {
		if _, ok := m[key].(float64); ok {
			delete(m, key)
		}
	}
}

func assertMapSlice(t *testing.T, expected, actual []map[string]interface{}, processors ...processor) {
	t.Helper()
	comparedSlice := []map[string]interface{}{}
	for _, a := range actual {
		compared := map[string]interface{}{}
		for k, v := range a {
			compared[k] = v
		}
		for _, p := range processors {
			p(t, compared)
		}
		comparedSlice = append(comparedSlice, compared)
	}
	assert.Equal(t, expected, comparedSlice)
}
