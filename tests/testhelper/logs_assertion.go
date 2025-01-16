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

func (a *LogsAssertion) Assert(t *testing.T, expected []map[string]interface{}) {
	actual := parseJsonLines(t, a.buf.Bytes())
	assertMapSlice(t, expected, actual, "time")
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

func assertMapSlice(t *testing.T, expected, actual []map[string]interface{}, ignoredFields ...string) {
	t.Helper()
	wellFormedActual := []map[string]interface{}{}
	for _, a := range actual {
		wellFormed := map[string]interface{}{}
		for k, v := range a {
			if contains(ignoredFields, k) {
				continue
			}
			wellFormed[k] = v
		}
		wellFormedActual = append(wellFormedActual, wellFormed)
	}
	assert.Equal(t, expected, wellFormedActual)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
