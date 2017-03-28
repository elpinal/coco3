package editor

import (
	"context"
	"io"
	"unicode/utf8"
)

type RuneAddReader struct {
	ird io.RuneReader
	s   []rune
	n   int
	ctx context.Context
}

func NewReaderContext(ctx context.Context, rd io.RuneReader) *RuneAddReader {
	return &RuneAddReader{ird: rd, ctx: ctx}
}

func NewReader(rd io.RuneReader) *RuneAddReader {
	return &RuneAddReader{ird: rd, ctx: context.Background()}
}

type runeRead struct {
	r    rune
	size int
	err  error
}

func (rd *RuneAddReader) readRune() chan runeRead {
	ch := make(chan runeRead)
	go func() {
		r, size, err := rd.ird.ReadRune()
		ch <- runeRead{r, size, err}
	}()
	return ch
}

func (rd *RuneAddReader) ReadRune() (r rune, size int, err error) {
	if len(rd.s) <= rd.n {
		ch := rd.readRune()
		defer func() { close(ch) }()
		select {
		case rr := <-ch:
			return rr.r, rr.size, rr.err
		case <-rd.ctx.Done():
			return 0, 0, io.EOF
		}
	}
	select {
	case <-rd.ctx.Done():
		return 0, 0, io.EOF
	default:
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
