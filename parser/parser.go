// Package parser implements a sam address parser
package parser

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/tudurom/usam/lex"
)

type Address struct {
	Simple SimpleAddress
	Left   *Address
	Next   *Address
}

type SimpleAddress interface {
	SimpleAddress() string // for differentiation
}

type CharAddress int

func (*CharAddress) SimpleAddress() string { return "#" }

type LineAddress int

func (LineAddress) SimpleAddress() string { return "1" }

type RegexAddress string

func (RegexAddress) SimpleAddress() string { return "/" }

type BackwardsRegexAddress string

func (BackwardsRegexAddress) SimpleAddress() string { return "?" }

type DotAddress bool

func (DotAddress) SimpleAddress() string { return "." }

type DollarAddress bool

func (DollarAddress) SimpleAddress() string { return "$" }

type PlusAddress bool

func (PlusAddress) SimpleAddress() string { return "+" }

type MinusAddress bool

func (MinusAddress) SimpleAddress() string { return "-" }

type CommaAddress bool

func (CommaAddress) SimpleAddress() string { return "," }

type SemicolonAddress bool

func (SemicolonAddress) SimpleAddress() string { return ";" }

type Parser struct {
	s         *lex.Scanner
	prevTok   lex.Token
	prevValue string
	hasPrev   bool
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: lex.NewScanner(r)}
}

func ParseString(s string) (*Address, error) {
	r := strings.NewReader(s)
	p := NewParser(r)
	return p.Parse()
}

func IsLowPrecedence(sa SimpleAddress) bool {
	_, isComma := sa.(*CommaAddress)
	_, isSemicolon := sa.(*SemicolonAddress)
	return isComma || isSemicolon
}

func IsHighPrecedence(sa SimpleAddress) bool {
	_, isPlus := sa.(*PlusAddress)
	_, isMinus := sa.(*MinusAddress)
	return isPlus || isMinus
}

func (p *Parser) scan() (tok lex.Token, value string) {
	if p.hasPrev {
		p.hasPrev = false
		return p.prevTok, p.prevValue
	}
	tok, value = p.s.Scan()
	p.prevTok = tok
	p.prevValue = value
	return
}

func (p *Parser) scanIgnoreWS() (lex.Token, string) {
	tok, val := p.scan()
	if tok == lex.WS {
		tok, val = p.scan()
	}

	return tok, val
}

func (p *Parser) unscan() {
	p.hasPrev = true
}

func (p *Parser) simpleAddress() *Address {
	addr := &Address{}
	tok, val := p.scanIgnoreWS()
	switch tok {
	case lex.CharAddr:
		i, err := strconv.Atoi(val[1:])
		if val[1:] != "" {
			if err != nil {
				panic(err)
			}
		} else {
			i = 1
		}
		ca := CharAddress(i)
		addr.Simple = &ca
	case lex.LineAddr:
		i, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		la := LineAddress(i)
		addr.Simple = &la
	case lex.Regexp:
		rx := RegexAddress(strings.TrimSuffix(strings.TrimPrefix(val, "/"), "/"))
		addr.Simple = &rx
	case lex.BackwardsRegexp:
		rx := BackwardsRegexAddress(strings.TrimSuffix(strings.TrimPrefix(val, "?"), "?"))
		addr.Simple = &rx
	case lex.Dot:
		da := DotAddress(true)
		addr.Simple = &da
	case lex.Dollar:
		da := DollarAddress(true)
		addr.Simple = &da
	case lex.Plus:
		pa := PlusAddress(true)
		addr.Simple = &pa
	case lex.Minus:
		ma := MinusAddress(true)
		addr.Simple = &ma
	default:
		p.unscan()
		return nil
	}

	if addr.Next = p.simpleAddress(); addr.Next != nil &&
		!IsHighPrecedence(addr.Next.Simple) &&
		!IsHighPrecedence(addr.Simple) {
		pa := PlusAddress(true)
		x := &Address{
			Simple: &pa,
			Next:   addr.Next,
		}
		addr.Next = x
	}
	//addr.Next = p.simpleAddress()

	return addr
}

// FillDefaults fills in defaults for compound addresses (+/-, ,/;)
func FillDefaults(addr *Address) *Address {
	if addr == nil {
		return nil
	}
	cur := addr
	var prev *Address
	for cur != nil {
		if IsHighPrecedence(cur.Simple) {
			/*
				A high precedence compound is of the form
					a1+a2
				or
					a1-a2
			*/

			// if a1 is missing, we put the dot
			if prev == nil {
				d := DotAddress(true)
				a := &Address{
					Simple: &d,
					Next:   addr,
				}
				addr = a
			}
			// if a2 is missing, we put the address to one line
			// so it will either add a line, or subtract a line
			if cur.Next == nil || IsHighPrecedence(cur.Next.Simple) {
				l := LineAddress(1)
				a := &Address{
					Simple: &l,
					Next:   cur.Next,
				}
				cur.Next = a
			}
		} else if IsLowPrecedence(cur.Simple) {
			/*
				A low precedence compound is of the form
					a1,a2
				or
					a1;a2
			*/

			cur.Left = FillDefaults(cur.Left)
			// if a1 is missing, we put the null line
			if cur.Left == nil {
				z := LineAddress(0)
				cur.Left = &Address{Simple: &z}
			}
			// if a2 is missing, we put the end of the file (dollar)
			if cur.Next == nil || IsLowPrecedence(cur.Next.Simple) {
				d := DollarAddress(true)
				cur.Next = &Address{
					Simple: &d,
					Next:   cur.Next,
				}
			}
		}
		prev = cur
		cur = cur.Next
	}
	return addr
}

func (p *Parser) Parse() (*Address, error) {
	addr := &Address{}

	addr.Left = p.simpleAddress()
	tok, val := p.scanIgnoreWS()
	if tok == lex.Comma {
		ca := CommaAddress(true)
		addr.Simple = &ca
	} else if tok == lex.Semicolon {
		sa := SemicolonAddress(true)
		addr.Simple = &sa
	} else if tok == lex.Illegal {
		return nil, errors.New("Illegal token '" + val + "'")
	} else {
		// golint will say to drop the else, don't do it
		return FillDefaults(addr.Left), nil
	}

	next, err := p.Parse()
	if err != nil {
		return nil, err
	}
	addr.Next = next

	if next != nil && IsLowPrecedence(next.Simple) && next.Left == nil {
		return nil, errors.New("Eaddress")
	}
	return FillDefaults(addr), nil
}
