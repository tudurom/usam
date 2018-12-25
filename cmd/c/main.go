package main

import (
	"fmt"
	"os"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/cliutil"
	"github.com/tudurom/usam/cliutil/pipeformat"
)

func usage() {
	fmt.Println("Usage: c <text>")
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	pf, err := pipeformat.Process()
	if err != nil {
		cliutil.Err(err)
	}

	fmt.Println(pf.Filename)

	a := usam.Address{Buffer: pf.Buffer, R: pf.Buffer.Dot}
	a, err = usam.ResolveAddress(pf.Addresses[0], a, 0)
	if err != nil {
		cliutil.Err(err)
	}
	pf.Buffer.Data = append(pf.Buffer.Data[:a.R.P1], append([]byte(os.Args[1]), pf.Buffer.Data[a.R.P2:]...)...)
	err = pf.Buffer.Save(pf.Filename)
	if err != nil {
		cliutil.Err(err)
	}
	fmt.Printf("#%d,#%d\n", a.R.P1, a.R.P1+len(os.Args[1]))
}
