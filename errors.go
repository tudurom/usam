package usam

import "errors"

// ErrOutOfRange is thrown when an address is outside the buffer's range
var ErrOutOfRange = errors.New("out of range")

// ErrWrongOrder is thrown when compound addresses are wrongly ordered
var ErrWrongOrder = errors.New("wrong order")

// ErrNoMatch is thrown when there is no match for a regexp
var ErrNoMatch = errors.New("no match")
