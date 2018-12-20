package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func usage() {
	fmt.Println("Usage: e <filename>")
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	_, err := os.Stat(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	abs, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(abs)
	fmt.Println("0")
}
