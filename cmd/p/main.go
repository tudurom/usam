package main

import (
	"fmt"
	"os"

	"github.com/tudurom/usam/parser"

	"github.com/tudurom/usam/cliutil"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/cliutil/pipeformat"
)

func usage() {
	fmt.Println("Usage: p [dot]")
}

func main() {
	cliutil.Name = "p"
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	pf, err := pipeformat.Process()
	if err != nil {
		cliutil.Err(err)
	}

	rarg := "."
	if len(os.Args) == 2 {
		rarg = os.Args[1]
	}
	aarg, err := parser.ParseString(rarg)
	if err != nil {
		cliutil.Err(err)
	}

	for _, ap := range pf.Addresses {
		a, err := usam.ResolveAddress(pf.Buffer.NewAddress(), ap, aarg)
		if err != nil {
			cliutil.Err(err)
		}
		fmt.Print(string(pf.Buffer.Data[a.R.P1:a.R.P2]))
	}
}
