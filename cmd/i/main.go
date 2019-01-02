package main

import (
	"os"

	"github.com/tudurom/usam/cliutil"
)

func main() {
	cliutil.C("i", func(text, dot []byte) []byte {
		r := make([]byte, len(text)+len(dot))
		copy(r, text)
		copy(r[len(text):], dot)
		return r
	}, os.Args)
}
