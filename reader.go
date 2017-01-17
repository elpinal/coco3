package main

import (
	"bufio"
	"io"
)

type Reader struct {
	src      []rune
	rd       *bufio.Reader
	offset   int
	lineHead int
	line     int
}

func NewReader(rd io.Reader) *Reader {
	r, ok := rd.(*bufio.Reader)
	if !ok {
		r = bufio.NewReader(rd)
	}
	return &Reader{rd: r}
}

func (r *Reader) read() (rune, error) {
	r.offset++
	ch, _, err := r.rd.ReadRune()
	if err != nil {
		return -1, err
	}
	return ch, nil
}

func (r *Reader) Read() (rune, error) {
	ch, err := r.read()
	return ch, err
}
