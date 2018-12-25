package pipeformat

import (
	"bufio"
	"io"
	"os"

	"github.com/tudurom/usam"
	"github.com/tudurom/usam/parser"
)

type PipeFormat struct {
	Filename  string
	Buffer    *usam.Buffer
	Addresses []*parser.Address
}

func process(r io.Reader) (PipeFormat, error) {
	pf := PipeFormat{}
	s := bufio.NewScanner(r)
	if s.Scan() {
		pf.Filename = s.Text()
	}
	if err := s.Err(); err != nil {
		return PipeFormat{}, err
	}

	var err error
	pf.Buffer, err = usam.NewBufferFromFile(pf.Filename)
	if err != nil {
		return PipeFormat{}, err
	}

	for s.Scan() {
		a, err := parser.ParseString(s.Text())
		if err != nil {
			// the pipe format is fully automated
			// no error is accepted
			panic(err)
		}
		pf.Addresses = append(pf.Addresses, a)
	}
	if err := s.Err(); err != nil {
		return PipeFormat{}, err
	}

	return pf, nil
}

func Process() (PipeFormat, error) {
	return process(os.Stdin)
}
