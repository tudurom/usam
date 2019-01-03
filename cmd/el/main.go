package main

import (
	"fmt"
	"os"

	"github.com/tudurom/usam/cliutil"

	"github.com/tudurom/usam/cliutil/pipeformat"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/parser"
)

func usage() {
	fmt.Println("Usage: el <address>")
}

func main() {
	cliutil.Name = "el"
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	pf, err := pipeformat.Process()
	if err != nil {
		cliutil.Err(err)
	}

	a1 := pf.Addresses[len(pf.Addresses)-1]
	a2, err := parser.ParseString(os.Args[1])
	if err != nil {
		cliutil.Err(err)
	}

	a, err := usam.ResolveAddress(pf.Buffer.NewAddress(), a1)
	if err != nil {
		cliutil.Err(err)
	}
	a.Buffer.Dot = a.R
	a, err = usam.ResolveAddress(a, a2)
	if err != nil {
		cliutil.Err(err)
	}

	fmt.Println(pf.Filename)
	fmt.Printf("#%d,#%d\n", a.R.P1, a.R.P2)
}
