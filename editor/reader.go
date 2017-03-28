package editor

import (
	"io"
	"unicode/utf8"
)

type RuneAddReader struct {
	ird io.RuneReader
	s   []rune
	n   int
}

func NewReader(rd io.RuneReader) *RuneAddReader {
	return &RuneAddReader{ird: rd}
}

func (rd *RuneAddReader) ReadRune() (r rune, size int, err error) {
	if len(rd.s) <= rd.n {
		return rd.ird.ReadRune()
	}
	r = rd.s[rd.n]
	size = utf8.RuneLen(r)
	rd.n++
	return r, size, err
}

func (rd *RuneAddReader) Add(s []rune) {
	rd.s = append(rd.s[rd.n:], s...)
	rd.n = 0
}
