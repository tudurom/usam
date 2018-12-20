package usam

import (
	"io"
	"io/ioutil"
	"os"
)

type Buffer struct {
	Data []byte
	Dot  Range
}

func NewBuffer(r io.Reader) (*Buffer, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &Buffer{
		Data: data,
		Dot:  Range{0, 0},
	}, nil
}

func NewBufferFromFile(fn string) (*Buffer, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := NewBuffer(f)
	if err != nil {
		return nil, err
	}

	return b, nil
}
