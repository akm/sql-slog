package sqlslog

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"strconv"
	"testing"
)

func idGenAttemptsFromEnv(t *testing.T) int {
	t.Helper()
	idGenAttemptsStr := os.Getenv("ID_GEN_ATTEMPTS")
	if idGenAttemptsStr == "" {
		idGenAttemptsStr = "1000"
	}
	idGenAttempts, err := strconv.Atoi(idGenAttemptsStr)
	if err != nil {
		t.Fatalf("strconv.Atoi: %v", err)
	}
	return idGenAttempts
}

func TestIDGeneratorDefault(t *testing.T) {
	t.Parallel()
	idGenAttempts := idGenAttemptsFromEnv(t)
	idGen := IDGeneratorDefault
	values := make([]string, idGenAttempts)
	for i := range idGenAttempts {
		values[i] = idGen()
	}
	for _, v := range values {
		if len(v) != defaultIDLength {
			t.Errorf("len(v) = %d, want %d", len(v), defaultIDLength)
		}
	}
	slices.Sort(values)
	compactValues := slices.Compact(values)
	if len(compactValues) < idGenAttempts {
		t.Errorf("len(compactValues) = %d, want %d", len(compactValues), idGenAttempts)
	}
}

func TestRandIntIDGenerator(t *testing.T) {
	t.Parallel()
	idGenAttempts := idGenAttemptsFromEnv(t)

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
			idGen := RandIntIDGenerator(
				rand.Int, // Use rand.Int from math/rand/v2
				defaultIDLetters,
				tc.length,
			)
			values := make([]string, idGenAttempts)
			for i := range idGenAttempts {
				values[i] = idGen()
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

func TestRandReadGenerator(t *testing.T) {
	t.Run("valid case", func(t *testing.T) {
		t.Parallel()
		idGenAttempts := idGenAttemptsFromEnv(t)

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
				idGen := RandReadGenerator(
					cryptorand.Read, // Use rand.Read from crypto/rand
					defaultIDLetters,
					tc.length,
				)
				values := make([]string, idGenAttempts)
				for i := range idGenAttempts {
					var err error
					values[i], err = idGen()
					if err != nil {
						t.Errorf("idGen: %v", err)
					}
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
	})

	t.Run("with error", func(t *testing.T) {
		idGen := RandReadGenerator(
			func([]byte) (int, error) { return 0, fmt.Errorf("unexpected error") },
			defaultIDLetters,
			defaultIDLength,
		)
		t.Run("return error", func(t *testing.T) {
			if _, err := idGen(); err == nil {
				t.Error("err = nil, want error")
			}
		})
		t.Run("suppress error", func(t *testing.T) {
			suppressedIDGen := IDGenErrorSuppressor(idGen,
				func(err error) string {
					return "suppressed"
				},
			)
			id := suppressedIDGen()
			if id != "suppressed" {
				t.Errorf("id = %q, want %q", id, "suppressed")
			}
		})
	})
}
