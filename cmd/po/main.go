package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tudurom/usam/cliutil"
)

func usage() {
	fmt.Println("Usage: po")
}

func main() {
	cliutil.Name = "po"
	if len(os.Args) != 1 {
		usage()
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		cliutil.Err(err)
	}

	f, err := ioutil.TempFile("", "usam-")
	if err != nil {
		cliutil.Err(err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		cliutil.Err(err)
	}

	fmt.Println(f.Name())
	fmt.Println("0")
}
