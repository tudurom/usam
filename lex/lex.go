// Package lex implements the lexer for the sam address parser
package lex

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// Token represents a symbolic token in the address to be parsed
type Token int

var eof = rune(0)

const (
	Illegal Token = iota
	EOF
	WS

	CharAddr
	LineAddr
	Regexp
	BackwardsRegexp

	Dot
	Plus
	Minus
	Comma
	Semicolon
	Dollar
)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() {
	s.r.UnreadRune()
}

func (s *Scanner) scanWS() (tok Token, value string) {
	var buf bytes.Buffer
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !unicode.IsSpace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *Scanner) scanCharAddress() (tok Token, value string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read()) // #

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !unicode.IsDigit(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return CharAddr, buf.String()
}

func (s *Scanner) scanLineAddress() (tok Token, value string) {
	var buf bytes.Buffer

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !unicode.IsDigit(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return LineAddr, buf.String()
}

func (s *Scanner) scanRegexp(reverse bool) (tok Token, value string) {
	var buf bytes.Buffer
	closed := false
	escaping := false
	char := '/'
	if reverse {
		char = '?'
	}
	buf.WriteRune(s.read())
	for {
		ch := s.read()
		fmt.Println("Scanning", ch)
		if ch == eof {
			break
		} else if ch == char {
			if escaping {
				buf.WriteRune(ch)
				escaping = false
			} else {
				closed = true
				break
			}
		} else {
			if escaping {
				buf.WriteRune(ch)
				escaping = false
			} else if ch == '\\' {
				escaping = true
			} else {
				buf.WriteRune(ch)
			}
		}
	}

	if closed {
		if reverse {
			return BackwardsRegexp, buf.String()
		}
		return Regexp, buf.String()
	}

	return Illegal, strings.TrimSpace(buf.String())
}

func (s *Scanner) Scan() (tok Token, value string) {
	ch := s.read()

	if unicode.IsSpace(ch) {
		s.unread()
		return s.scanWS()
	}

	switch {
	case ch == eof:
		return EOF, ""
	case ch == '#':
		s.unread()
		return s.scanCharAddress()
	case unicode.IsDigit(ch):
		s.unread()
		return s.scanLineAddress()
	case ch == '/':
		s.unread()
		return s.scanRegexp(false)
	case ch == '?':
		s.unread()
		return s.scanRegexp(true)

	case ch == '.':
		return Dot, "."
	case ch == '+':
		return Plus, "+"
	case ch == '-':
		return Minus, "-"
	case ch == ',':
		return Comma, ","
	case ch == ';':
		return Semicolon, ";"
	case ch == '$':
		return Dollar, "$"
	}

	return Illegal, string(ch)
}
