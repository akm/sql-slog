package sqlslog

import (
	"fmt"
	"testing"
)

func TestWrapConn(t *testing.T) {
	if wrapConn(nil, nil) != nil {
		t.Fatal("Expected nil")
	}
}

func TestConnExecContextErrorHandler(t *testing.T) {
	errHandler := ConnExecContextErrorHandler("mysql")
	complete, attrs := errHandler(fmt.Errorf("dummy"))
	if complete {
		t.Fatal("Expected false")
	}
	if attrs != nil {
		t.Fatal("Expected nil")
	}
}

func TestConnQueryContextErrorHandler(t *testing.T) {
	errHandler := ConnQueryContextErrorHandler("mysql")
	complete, attrs := errHandler(fmt.Errorf("dummy"))
	if complete {
		t.Fatal("Expected false")
	}
	if attrs != nil {
		t.Fatal("Expected nil")
	}
}
