package cliutil

import (
	"fmt"
	"os"
)

func Err(err error) {
	fmt.Fprintln(os.Stderr, os.Args[0]+":", err)
	os.Exit(1)
}
