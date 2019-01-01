package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/tudurom/usam/parser"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/cliutil"
	"github.com/tudurom/usam/cliutil/pipeformat"
)

func usage() {
	fmt.Println("Usage: a <dot> <text>")
}

func main() {
	cliutil.Name = "a"
	if len(os.Args) != 3 {
		usage()
		os.Exit(1)
	}

	pf, err := pipeformat.Process()
	if err != nil {
		cliutil.Err(err)
	}

	text := []byte(os.Args[2])
	aarg, err := parser.ParseString(os.Args[1])
	if err != nil {
		cliutil.Err(err)
	}

	var as []usam.Address
	for _, ap := range pf.Addresses {
		a, err := usam.ResolveAddress(pf.Buffer.NewAddress(), ap, aarg)
		if err != nil {
			cliutil.Err(err)
		}
		as = append(as, a)
	}
	sort.Sort(cliutil.ByP1(as))

	fmt.Println(pf.Filename)
	i := 0
	for _, a := range as {
		pf.Buffer.Data = append(pf.Buffer.Data[:a.R.P2+i], append(text, pf.Buffer.Data[i+a.R.P2:]...)...)
		fmt.Printf("#%d,#%d\n", a.R.P2+i, a.R.P2+i+len(text))
		i += len(text)
	}
	err = pf.Buffer.Save(pf.Filename)
	if err != nil {
		cliutil.Err(err)
	}
}
