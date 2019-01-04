package main

import (
	"os"

	"github.com/tudurom/usam/cliutil"
)

func main() {
	cliutil.C("d", func(text, dot []byte) []byte {
		return nil
	}, os.Args)
}
