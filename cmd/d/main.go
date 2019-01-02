package main

import (
	"github.com/tudurom/usam/cliutil"
	"os"
)

func main() {
	cliutil.C("d", func(text, dot []byte) []byte {
		return nil
	}, os.Args)
}
