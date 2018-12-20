package main

import (
	"fmt"
	"os"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/parser"
)

func usage() {
	fmt.Println("Usage: ca <address>")
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	var fn string
	fmt.Scanln(&fn)
	b, err := usam.NewBufferFromFile(fn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var a1r string
	fmt.Scanln(&a1r)
	a1, err := parser.ParseString(a1r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	a2, err := parser.ParseString(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	a := usam.Address{Buffer: b, R: b.Dot}
	a, err = usam.ResolveAddress(a1, a, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	a.Buffer.Dot = a.R
	a, err = usam.ResolveAddress(a2, a, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(fn)
	fmt.Printf("#%d,#%d\n", a.R.P1, a.R.P2)
}
