package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tudurom/usam/cliutil"
)

func usage() {
	fmt.Println("Usage: e <filename>")
}

func main() {
	cliutil.Name = "e"
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	_, err := os.Stat(os.Args[1])
	if err != nil {
		cliutil.Err(err)
	}

	abs, err := filepath.Abs(os.Args[1])
	if err != nil {
		cliutil.Err(err)
	}

	fmt.Println(abs)
	fmt.Println("0")
}
