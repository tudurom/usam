package cliutil

import (
	"fmt"
	"os"

	"github.com/tudurom/usam"
)

// ByP1 implements sort.Interface for []usam.Address
// based on the left side of the range (.R.P1)
type ByP1 []usam.Address

func (a ByP1) Len() int           { return len(a) }
func (a ByP1) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByP1) Less(i, j int) bool { return a[i].R.P1 < a[j].R.P1 }

var Name string

func Err(err error) {
	fmt.Fprintln(os.Stderr, Name+":", err)
	os.Exit(1)
}
