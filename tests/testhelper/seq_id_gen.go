package testhelper

import "fmt"

func NewSeqIdGenerator(start int) func() string {
	seq := start
	return func() string {
		seq++
		return fmt.Sprintf("%04d", seq)
	}
}
