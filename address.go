package usam

import (
	"regexp"

	"github.com/tudurom/usam/parser"
)

// Address is a chunk of consecutive text (a range) in a buffer
type Address struct {
	R      Range
	Buffer *Buffer
}

func regexAddress(pattern string, addr Address, sign int) (Address, error) {
	// location of the regex. the new dot
	var loc []int
	// the reference point
	// it's p2, the rightmost end of the dot if we're going forward
	// or p1 if backwards
	var l int
	// the regex must be in multiline mode
	re, err := regexp.Compile("(?m)" + pattern)
	if err != nil {
		return Address{}, err
	}
	if sign >= 0 {
		l = addr.R.P2
		if loc = re.FindIndex(addr.Buffer.Data[l:]); loc == nil {
			return Address{}, ErrNoMatch
		}
		loc[0] += l
		loc[1] += l

		if loc[0] == loc[1] && loc[0] == l {
			l++
			if l > len(addr.Buffer.Data) {
				l = 0
			}
			if loc = re.FindIndex(addr.Buffer.Data[l:]); loc == nil {
				panic("address")
			}
			loc[0] += l
			loc[1] += l
		}
	} else {
		l = addr.R.P1
		locs := re.FindAllIndex(addr.Buffer.Data[:l], -1)
		if locs == nil {
			return Address{}, ErrNoMatch
		}
		loc = locs[len(locs)-1]
		if loc[0] == loc[1] && loc[0] == l {
			l--
			if l < 0 {
				l = len(addr.Buffer.Data)
			}
			locs := re.FindAllIndex(addr.Buffer.Data[:l], -1)
			if locs == nil {
				panic("address")
			}
			loc = locs[len(locs)-1]
		}
	}

	return Address{
		R:      Range{loc[0], loc[1]},
		Buffer: addr.Buffer,
	}, nil
}

func lineAddress(la *parser.LineAddress, addr Address, sign int) (Address, error) {
	var a Address
	a.Buffer = addr.Buffer
	p := 0
	n := 0

	l := int(*la)
	if sign >= 0 {
		if l == 0 {
			// we are either talking about absolute line numbers
			// so l = 0 means the null line
			// or the address is already there
			// so we return the null line
			if sign == 0 || addr.R.P2 == 0 {
				a.R.P1 = 0
				a.R.P2 = 0
				return a, nil
			}
			// otherwise, we are looking for the end of line
			// you will see this below
			a.R.P1 = addr.R.P2
			p = addr.R.P2 - 1
		} else {
			// again, if it's the null line...
			if sign == 0 || addr.R.P2 == 0 {
				p = 0
				// n = 1 means we will skip the line we are on
				// n is the iterator
				n = 1
			} else {
				p = addr.R.P2 - 1
				// if we are just at the start of the line, we skip it too
				if addr.Buffer.Data[p] == '\n' {
					n = 1
				} else {
					n = 0
				}
				p++
			}
			// we count the lines now
			for n < l {
				if p >= len(addr.Buffer.Data) {
					return Address{}, ErrOutOfRange
				}
				if addr.Buffer.Data[p] == '\n' {
					n++
				}
				p++
			}
			// start of the line
			a.R.P1 = p
		}
		// find the end of the line
		for p < len(addr.Buffer.Data) && addr.Buffer.Data[p] != '\n' {
			p++
		}
		// end of the line
		a.R.P2 = p
		if p < len(addr.Buffer.Data) {
			a.R.P2++
		}
	} else {
		p = addr.R.P1
		if l == 0 {
			// we are looking for the 0th line,
			// relative from where we are now, backwards
			// so this means we are looking for the end of the previous line
			a.R.P2 = addr.R.P1
		} else {
			n = 0
			for n < l {
				if p == 0 {
					n++
					if n != l {
						// we are at the start of the buffer and
						// the search is not over. it's clearly an error.
						return Address{}, ErrOutOfRange
					}
				} else {
					c := addr.Buffer.Data[p-1]
					if c != '\n' || n+1 != l {
						p--
					}
					if c == '\n' {
						n++
					}
				}
			}
			a.R.P2 = p
			if p > 0 {
				p--
			}
		}
		// lines start after a newline
		for p > 0 && addr.Buffer.Data[p-1] != '\n' {
			p--
		}
		a.R.P1 = p
	}
	return a, nil
}

func charAddress(ca *parser.CharAddress, addr Address, sign int) (Address, error) {
	l := int(*ca)
	if sign == 0 {
		addr.R.P1 = l
		addr.R.P2 = l
	} else if sign < 0 {
		addr.R.P1 -= l
		addr.R.P2 -= l
	} else {
		addr.R.P1 += l
		addr.R.P2 += l
	}
	if addr.R.P1 < 0 || addr.R.P2 > len(addr.Buffer.Data) {
		return Address{}, ErrOutOfRange
	}
	return addr, nil
}

func ResolveAddress(base Address, aps ...*parser.Address) (Address, error) {
	var err error
	for _, ap := range aps {
		base, err = resolveAddress(ap, base, 0)
		if err != nil {
			return Address{}, err
		}
		base.Buffer.Dot = base.R
	}
	return base, nil
}

func resolveAddress(ap *parser.Address, a Address, sign int) (Address, error) {
	for ap != nil {
		var err error
		switch ap.Simple.SimpleAddress() {
		case "1":
			a, err = lineAddress(ap.Simple.(*parser.LineAddress), a, sign)
		case "#":
			a, err = charAddress(ap.Simple.(*parser.CharAddress), a, sign)
		case "$":
			a.R.P1 = len(a.Buffer.Data)
			a.R.P2 = len(a.Buffer.Data)
		case ".":
			a.R = a.Buffer.Dot
		case "?":
			sign = -sign
			if sign == 0 {
				sign = -1
			}
			a, err = regexAddress(string(*ap.Simple.(*parser.BackwardsRegexAddress)), a, sign)
		case "/":
			a, err = regexAddress(string(*ap.Simple.(*parser.RegexAddress)), a, sign)
		case ",":
			fallthrough
		case ";":
			var a1, a2 Address
			if ap.Left != nil {
				a1, err = resolveAddress(ap.Left, a, sign)
				if err != nil {
					return Address{}, err
				}
			} else {
				a1.Buffer = a.Buffer
				a1.R.P1 = 0
				a1.R.P2 = 0
			}
			if ap.Next != nil {
				a2, err = resolveAddress(ap.Next, a, sign)
				if err != nil {
					return Address{}, err
				}
			} else {
				a2.Buffer = a.Buffer
				a2.R.P1 = len(a.Buffer.Data)
				a2.R.P2 = len(a.Buffer.Data)
			}
			a.Buffer = a1.Buffer
			a.R.P1 = a1.R.P1
			a.R.P2 = a2.R.P2
			if a.R.P2 < a.R.P1 {
				return Address{}, ErrWrongOrder
			}
			return a, nil
		case "+":
			sign = +1
		case "-":
			sign = -1
		}

		if err != nil {
			return Address{}, err
		}
		ap = ap.Next
	}

	return a, nil
}
