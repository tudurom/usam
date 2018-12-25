package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/cliutil"
	"github.com/tudurom/usam/cliutil/pipeformat"
)

func usage() {
	fmt.Println("Usage: s <regexp> <text> [n|g]")
}

func main() {
	if len(os.Args) < 3 || len(os.Args) > 4 {
		usage()
		os.Exit(1)
	}

	n := 1
	if len(os.Args) == 4 {
		if os.Args[3] == "g" {
			n = -1
		} else {
			var err error
			n, err = strconv.Atoi(os.Args[3])
			if err != nil {
				cliutil.Err(err)
			}
		}
	}
	re, err := regexp.Compile(os.Args[1])
	if err != nil {
		cliutil.Err(err)
	}

	pf, err := pipeformat.Process()
	if err != nil {
		cliutil.Err(err)
	}

	a, err := usam.ResolveAddress(pf.Addresses[0], pf.Buffer.NewAddress(), 0)
	if err != nil {
		cliutil.Err(err)
	}

	tmpl := []byte(os.Args[2])
	if n == -1 {
		pf.Buffer.Data = append(
			pf.Buffer.Data[:a.R.P1],
			append(
				re.ReplaceAll(pf.Buffer.Data[a.R.P1:a.R.P2], tmpl),
				pf.Buffer.Data[a.R.P2:]...,
			)...)
	} else {
		index := re.FindAllSubmatchIndex(pf.Buffer.Data[a.R.P1:a.R.P2], n)[n-1]
		var result []byte
		result = re.Expand(result, tmpl, pf.Buffer.Data[a.R.P1:a.R.P2], index)
		pf.Buffer.Data = append(
			pf.Buffer.Data[:a.R.P1+index[0]],
			append(result, pf.Buffer.Data[a.R.P1+index[1]:]...)...)
	}
	if err = pf.Buffer.Save(pf.Filename); err != nil {
		cliutil.Err(err)
	}
}
