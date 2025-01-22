//go:build go1.23

package sqlslog

import (
	cryptoRand "crypto/rand"
	"fmt"
	"math/rand/v2"
)

type ChaCha8IDGenerator struct {
	*rand.ChaCha8
	length int
}

func NewChaCha8IDGenerator(length int) *ChaCha8IDGenerator {
	var s [32]byte
	if _, err := cryptoRand.Read(s[:]); err != nil {
		panic(fmt.Sprintf("sqldblogger: could not get random bytes from cryto/rand: '%s'", err.Error()))
	}

	return &ChaCha8IDGenerator{
		ChaCha8: rand.NewChaCha8(s),
		length:  length,
	}
}

const genratorChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"

func (g *ChaCha8IDGenerator) Generate() string {
	random := make([]byte, g.length)
	uid := make([]byte, g.length)

	g.ChaCha8.MarshalBinary()

	// rand.ChaCha8.Read is available since go1.23
	if _, err := g.Read(random); err != nil {
		panic(fmt.Sprintf("sql-slog: random read error from math/rand/v2 ChaCha8: %q", err.Error()))
	}
	for i := 0; i < g.length; i++ {
		uid[i] = genratorChars[random[i]&62]
	}
	return string(uid)
}
