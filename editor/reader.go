package editor

import (
	"context"
	"io"
)

type RecordableRuneReader struct {
	rd  io.RuneReader
	s   []rune
	ctx context.Context

	record bool
}

func NewReaderContext(ctx context.Context, rd io.RuneReader) *RecordableRuneReader {
	return &RecordableRuneReader{rd: rd, ctx: ctx}
}

func NewReader(rd io.RuneReader) *RecordableRuneReader {
	return NewReaderContext(context.Background(), rd)
}

type runeRead struct {
	r    rune
	size int
	err  error
}

func (rd *RecordableRuneReader) readRune() chan runeRead {
	ch := make(chan runeRead)
	go func() {
		r, size, err := rd.rd.ReadRune()
		ch <- runeRead{r, size, err}
	}()
	return ch
}

func (rd *RecordableRuneReader) ReadRune() (r rune, size int, err error) {
	ch := rd.readRune()
	defer func() { close(ch) }()
	select {
	case rr := <-ch:
		if rd.record {
			rd.s = append(rd.s, rr.r)
		}
		return rr.r, rr.size, rr.err
	case <-rd.ctx.Done():
		return 0, 0, io.EOF
	}
}

func (rd *RecordableRuneReader) Record() {
	rd.record = true
}

func (rd *RecordableRuneReader) Stop() []rune {
	rd.record = false
	ret := rd.s
	rd.s = nil
	return ret
}
