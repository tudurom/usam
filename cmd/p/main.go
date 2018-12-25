package main

import (
	"fmt"
	"os"

	"github.com/tudurom/usam/cliutil"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/cliutil/pipeformat"
)

func main() {
	pf, err := pipeformat.Process()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, ap := range pf.Addresses {
		a := usam.Address{R: pf.Buffer.Dot, Buffer: pf.Buffer}
		a, err = usam.ResolveAddress(ap, a, 0)
		if err != nil {
			cliutil.Err(err)
		}
		fmt.Print(string(pf.Buffer.Data[a.R.P1:a.R.P2]))
	}
}
