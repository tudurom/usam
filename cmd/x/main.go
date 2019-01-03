package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/tudurom/usam"

	"github.com/tudurom/usam/cliutil"

	"github.com/tudurom/usam/cliutil/pipeformat"
)

func usage() {
	fmt.Println("Usage: x <regex>")
}

func main() {
	cliutil.Name = "x"
	if len(os.Args) > 2 {
		usage()
		os.Exit(1)
	}

	var re *regexp.Regexp
	if len(os.Args) == 1 {
		re = regexp.MustCompile("(?m).*\n")
	} else {
		var err error
		re, err = regexp.Compile("(?m)" + os.Args[1])
		if err != nil {
			cliutil.Err(err)
		}
	}

	pf, err := pipeformat.Process()
	if err != nil {
		cliutil.Err(err)
	}

	fmt.Println(pf.Filename)
	for _, ap := range pf.Addresses {
		a, err := usam.ResolveAddress(pf.Buffer.NewAddress(), ap)
		if err != nil {
			cliutil.Err(err)
		}

		matches := re.FindAllIndex(pf.Buffer.Data[a.R.P1:a.R.P2], -1)
		for _, match := range matches {
			fmt.Printf("#%d,#%d\n", a.R.P1+match[0], a.R.P1+match[1])
		}
	}
}
