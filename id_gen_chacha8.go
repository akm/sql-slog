package sqlslog

import (
	cryptoRand "crypto/rand"
	"fmt"
	"math/rand/v2"
)

// ChaCha8IDGenerator generates random IDs using ChaCha8 in math/rand/v2.
// The seed is generated using crypto/rand.
type ChaCha8IDGenerator struct {
	*rand.Rand
	length int
}

// NewChaCha8IDGenerator creates a new ChaCha8IDGenerator with a random seed.
func NewChaCha8IDGenerator(length int) *ChaCha8IDGenerator {
	var s [32]byte
	if _, err := cryptoRand.Read(s[:]); err != nil {
		panic(fmt.Sprintf("sql-slog: could not get random bytes from crypto/rand: %q", err.Error()))
	}

	return &ChaCha8IDGenerator{
		Rand:   rand.New(rand.NewChaCha8(s)), // nolint:gosec
		length: length,
	}
}

var generatorRunes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")

// Generate generates a random ID.
func (g *ChaCha8IDGenerator) Generate() string {
	r := make([]byte, g.length)
	runesLen := len(generatorRunes)
	for i := range g.length {
		r[i] = generatorRunes[g.IntN(runesLen)]
	}
	return string(r)
}
