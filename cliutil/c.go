package cliutil

import (
	"fmt"
	"os"
	"sort"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/cliutil/pipeformat"
	"github.com/tudurom/usam/parser"
)

func cUsage() {
	if Name == "d" {
		fmt.Printf("Usage: %s [dot]\n", Name)
	} else {
		fmt.Printf("Usage: %s <text> [dot]\n", Name)
	}
}

type ChangeFunc func(text, dot []byte) []byte

func C(name string, cf ChangeFunc, args []string) {
	Name = name
	if (name != "d" && (len(args) < 2 || len(args) > 3)) ||
		(name == "d" && len(args) > 2) {
		cUsage()
		os.Exit(1)
	}

	pf, err := pipeformat.Process()
	if err != nil {
		Err(err)
	}

	var argtext []byte
	if name != "d" {
		argtext = []byte(args[1])
	}
	rarg := "."
	if len(args) == 3 {
		rarg = args[2]
	} else if name == "d" && len(args) == 2 {
		rarg = args[1]
	}
	aarg, err := parser.ParseString(rarg)
	if err != nil {
		Err(err)
	}

	var as []usam.Address
	for _, ap := range pf.Addresses {
		a, err := usam.ResolveAddress(pf.Buffer.NewAddress(), ap, aarg)
		if err != nil {
			Err(err)
		}
		as = append(as, a)
	}
	sort.Sort(ByP1(as))

	fmt.Println(pf.Filename)
	i := 0
	for _, a := range as {
		text := cf(argtext, pf.Buffer.Data[a.R.P1+i:a.R.P2+i])
		pf.Buffer.Data = append(pf.Buffer.Data[:a.R.P1+i], append(text, pf.Buffer.Data[i+a.R.P2:]...)...)
		pf.NewOutput(usam.Range{a.R.P1 + i, a.R.P1 + i + len(text)})
		i += len(text) - (a.R.P2 - a.R.P1)
	}
	err = pf.Buffer.Save(pf.Filename)
	if err != nil {
		Err(err)
	}
	pf.Print()
}
