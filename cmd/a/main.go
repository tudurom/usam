package main

import (
	"os"

	"github.com/tudurom/usam/cliutil"
)

func main() {
	cliutil.C("a", func(text, dot []byte) []byte {
		r := make([]byte, len(text)+len(dot))
		// I use copy because append updates the slice in-place
		copy(r, dot)
		copy(r[len(dot):], text)
		return r
	}, os.Args)
}
