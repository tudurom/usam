package main

import (
	"fmt"
	"os"

	"github.com/tudurom/usam/cliutil/pipeformat"

	"github.com/tudurom/usam/cliutil"
)

func usage() {
	fmt.Println("Usage: pc")
}

func main() {
	cliutil.Name = "pc"
	if len(os.Args) != 1 {
		usage()
		os.Exit(1)
	}

	pf, err := pipeformat.Process()
	if err != nil {
		cliutil.Err(err)
	}

	if _, err = os.Stdout.Write(pf.Buffer.Data); err != nil {
		cliutil.Err(err)
	}
	if err := os.Remove(pf.Filename); err != nil {
		cliutil.Err(err)
	}
}
