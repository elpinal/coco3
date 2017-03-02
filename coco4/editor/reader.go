package editor

import (
	"io"
	"sync"
	"unicode/utf8"
)

type RuneAddReader struct {
	ird io.RuneReader

	mu sync.Mutex
	ch chan rune
}

func NewReader(rd io.RuneReader) *RuneAddReader {
	return &RuneAddReader{ird: rd, ch: make(chan rune, 2)}
}

func (rd *RuneAddReader) ReadRune() (r rune, size int, err error) {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	select {
	case r = <-rd.ch:
		size = utf8.RuneLen(r)
		return r, size, err
	default:
		return rd.ird.ReadRune()
	}
}

func (rd *RuneAddReader) Add(s []rune) {
	rd.mu.Lock()
	go func() {
		for _, r := range s {
			rd.ch <- r
		}
		rd.mu.Unlock()
	}()
}
