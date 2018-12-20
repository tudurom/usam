package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/tudurom/usam/parser"
)

func main() {
	p := parser.NewParser(os.Stdin)
	x, err := p.Parse()
	if err != nil {
		panic(err)
	}
	spew.Printf("%#v\n", x)
}
