package usam

import (
	"io"
	"io/ioutil"
	"os"
)

// Buffer is a temporary representation of the file we are manipulating
type Buffer struct {
	Data []byte
	Dot  Range
}

// NewBuffer creates a buffer from a reader
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

// NewBufferFromFile creates a buffer from a file
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

// Save saves the buffer's contents to a file
func (b *Buffer) Save(fn string) error {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b.Data)
	return err
}

func (b *Buffer) NewAddress() Address {
	return Address{Buffer: b, R: b.Dot}
}
