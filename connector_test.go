package sqlslog

import (
	"fmt"
	"testing"
)

func TestConnectorConnectErrorHandler(t *testing.T) {
	testcases := []string{
		"mysql",
		"postgres",
	}
	for _, driverName := range testcases {
		t.Run(driverName, func(t *testing.T) {
			errHandler := ConnectorConnectErrorHandler(driverName)
			complete, attrs := errHandler(fmt.Errorf("dummy"))
			if complete {
				t.Fatal("Expected false")
			}
			if attrs != nil {
				t.Fatal("Expected nil")
			}
		})
	}
}
