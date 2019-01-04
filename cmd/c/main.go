package main

import (
	"os"

	"github.com/tudurom/usam/cliutil"
)

func main() {
	cliutil.C("c", func(text, dot []byte) []byte {
		return text
	}, os.Args)
}
