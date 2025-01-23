package testhelper

import "fmt"

type SeqIDGenerator struct {
	format string
	seq    int
}

func NewSeqIDGenerator() *SeqIDGenerator {
	return &SeqIDGenerator{
		format: "%04d",
		seq:    0,
	}
}

func (g *SeqIDGenerator) Generate() string {
	g.seq++
	return fmt.Sprintf(g.format, g.seq)
}

func (g *SeqIDGenerator) Next() string {
	return fmt.Sprintf(g.format, g.seq+1)
}

func (g *SeqIDGenerator) Set(v int) {
	g.seq = v
}
