package main

import (
	"fmt"
	"os"

	"github.com/tudurom/usam/cliutil"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/cliutil/pipeformat"
)

func usage() {
	fmt.Println("Usage: d")
}

func main() {
	cliutil.Name = "d"
	if len(os.Args) != 1 {
		usage()
		os.Exit(1)
	}

	pf, err := pipeformat.Process()
	if err != nil {
		cliutil.Err(err)
	}

	fmt.Println(pf.Filename)

	a, err := usam.ResolveAddress(pf.Buffer.NewAddress(), pf.Addresses[0])
	if err != nil {
		cliutil.Err(err)
	}
	pf.Buffer.Data = append(pf.Buffer.Data[:a.R.P1], pf.Buffer.Data[a.R.P2:]...)
	err = pf.Buffer.Save(pf.Filename)
	if err != nil {
		cliutil.Err(err)
	}
	fmt.Printf("#%d,#%d\n", a.R.P1, a.R.P1)
}
