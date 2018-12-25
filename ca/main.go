package main

import (
	"fmt"
	"os"

	"github.com/tudurom/usam/pipeformat"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/parser"
)

func usage() {
	fmt.Println("Usage: ca <address>")
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	pf, err := pipeformat.Process()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	a1 := pf.Addresses[0]
	a2, err := parser.ParseString(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	a := usam.Address{Buffer: pf.Buffer, R: pf.Buffer.Dot}
	a, err = usam.ResolveAddress(a1, a, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	a.Buffer.Dot = a.R
	a, err = usam.ResolveAddress(a2, a, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(pf.Filename)
	fmt.Printf("#%d,#%d\n", a.R.P1, a.R.P2)
}
