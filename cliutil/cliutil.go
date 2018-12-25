package cliutil

import (
	"fmt"
	"os"
)

func Err(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
