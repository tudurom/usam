package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/tudurom/usam/parser"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/cliutil"
	"github.com/tudurom/usam/cliutil/pipeformat"
)

func usage() {
	fmt.Println("Usage: s <regexp> <text> [n|g] [dot]")
}

func main() {
	cliutil.Name = "s"
	if len(os.Args) < 3 || len(os.Args) > 5 {
		usage()
		os.Exit(1)
	}

	n := 1
	all := false
	if len(os.Args) == 4 {
		if os.Args[3] == "g" {
			all = true
			n = -1
		} else {
			var err error
			n, err = strconv.Atoi(os.Args[3])
			if err != nil {
				cliutil.Err(err)
			}
		}
	}
	re, err := regexp.Compile("(?m)" + os.Args[1])
	if err != nil {
		cliutil.Err(err)
	}

	pf, err := pipeformat.Process()
	if err != nil {
		cliutil.Err(err)
	}

	rarg := "."
	if len(os.Args) == 5 {
		rarg = os.Args[4]
	}
	aarg, err := parser.ParseString(rarg)
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

	// XXX: it's ugly
	fmt.Println(pf.Filename)
	tmpl := []byte(os.Args[2])
	i := 0
	delta := 0
	for _, a := range as {
		submatches := re.FindAllSubmatchIndex(pf.Buffer.Data[i+a.R.P1:a.R.P2+i], n)
		if len(submatches) < n {
			cliutil.Err(usam.ErrNoMatch)
		}
		iterations := n
		if all {
			iterations = len(submatches)
		}
		for k := 0; k < iterations; k++ {
			if !all && k < iterations-1 {
				continue
			}

			index := submatches[k]
			var result []byte
			result = re.Expand(result, tmpl, pf.Buffer.Data[i+a.R.P1:a.R.P2+i+delta], index)
			for it := range index {
				index[it] += delta
			}
			pf.Buffer.Data = append(
				pf.Buffer.Data[:a.R.P1+i+index[0]],
				append(result, pf.Buffer.Data[a.R.P1+i+index[1]:]...)...)
			pf.NewOutput(usam.Range{P1: a.R.P1 + i + index[0], P2: a.R.P1 + i + index[1]})
			delta = len(result) - (index[1] - index[0])
			i += delta
		}
	}
	if err = pf.Buffer.Save(pf.Filename); err != nil {
		cliutil.Err(err)
	}
	pf.Print()
}
