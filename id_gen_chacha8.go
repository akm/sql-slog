package sqlslog

import (
	cryptoRand "crypto/rand"
	"fmt"
	"math/rand/v2"
)

type ChaCha8IDGenerator struct {
	*rand.Rand
	length int
}

func NewChaCha8IDGenerator(length int) *ChaCha8IDGenerator {
	var s [32]byte
	if _, err := cryptoRand.Read(s[:]); err != nil {
		panic(fmt.Sprintf("sql-slog: could not get random bytes from crypto/rand: %q", err.Error()))
	}

	return &ChaCha8IDGenerator{
		Rand:   rand.New(rand.NewChaCha8(s)),
		length: length,
	}
}

var generatorRunes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")

func (g *ChaCha8IDGenerator) Generate() string {
	r := make([]byte, g.length)
	runesLen := len(generatorRunes)
	for i := 0; i < g.length; i++ {
		r[i] = generatorRunes[g.IntN(runesLen)]
	}
	return string(r)
}
