package sqlslog

import (
	"testing"
)

func TestWrapConn(t *testing.T) {
	if wrapConn(nil, nil) != nil {
		t.Fatal("Expected nil")
	}
}
