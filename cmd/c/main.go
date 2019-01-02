package main

import (
	"github.com/tudurom/usam/cliutil"
	"os"
)


func main() {
	cliutil.C("c", func(text, dot []byte) []byte {
		return text
	}, os.Args)
}