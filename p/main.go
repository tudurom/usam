package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/parser"
)

func main() {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(string(buf)), "\n")

	b, err := usam.NewBufferFromFile(lines[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 1; i < len(lines); i++ {
		ap, err := parser.ParseString(lines[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		a := usam.Address{R: b.Dot, Buffer: b}
		a, err = usam.ResolveAddress(ap, a, 0)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(b.Data[a.R.P1:a.R.P2]))
	}
}
