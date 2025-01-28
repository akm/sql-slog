package public

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"testing"
)

func TestChaCha8Gen(t *testing.T) {
	t.Parallel()
	idGenAttemptsStr := os.Getenv("ID_GEN_ATTEMPTS")
	if idGenAttemptsStr == "" {
		idGenAttemptsStr = "1000"
	}
	idGenAttempts, err := strconv.Atoi(idGenAttemptsStr)
	if err != nil {
		t.Fatalf("strconv.Atoi: %v", err)
	}

	testCases := []struct {
		length int
	}{
		{length: 8},
		{length: 12},
		{length: 16},
		{length: 24},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("length %d", tc.length), func(t *testing.T) {
			t.Parallel()
			idGen := NewChaCha8IDGenerator(tc.length)
			values := make([]string, idGenAttempts)
			for i := range idGenAttempts {
				values[i] = idGen.Generate()
			}
			for _, v := range values {
				if len(v) != tc.length {
					t.Errorf("len(v) = %d, want %d", len(v), tc.length)
				}
			}
			slices.Sort(values)
			compactValues := slices.Compact(values)
			if len(compactValues) < idGenAttempts {
				t.Errorf("len(compactValues) = %d, want %d", len(compactValues), idGenAttempts)
			}
		})
	}
}
