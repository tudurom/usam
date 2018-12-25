package usam

import "fmt"

// Range is a simple, inclusive range
// [P1, P2]
type Range struct {
	P1 int
	P2 int
}

func (r Range) String() string {
	return fmt.Sprintf("{%d, %d}", r.P1, r.P2)
}
