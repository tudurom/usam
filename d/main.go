package main

import (
	"fmt"
	"os"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/pipeformat"
)

func usage() {
	fmt.Println("Usage: d")
}

func main() {
	if len(os.Args) != 1 {
		usage()
		os.Exit(1)
	}

	pf, err := pipeformat.Process()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(pf.Filename)

	a := usam.Address{Buffer: pf.Buffer, R: pf.Buffer.Dot}
	a, err = usam.ResolveAddress(pf.Addresses[0], a, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pf.Buffer.Data = append(pf.Buffer.Data[:a.R.P1], pf.Buffer.Data[a.R.P2:]...)
	err = pf.Buffer.Save(pf.Filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("#%d,#%d\n", a.R.P1, a.R.P1)
}
