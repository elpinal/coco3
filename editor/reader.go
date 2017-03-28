package editor

import (
	"io"
	"sync"
	"unicode/utf8"
)

type RuneAddReader struct {
	ird io.RuneReader
	ch  chan rune
	wg  *sync.WaitGroup
}

func NewReader(rd io.RuneReader) *RuneAddReader {
	return &RuneAddReader{ird: rd, ch: make(chan rune), wg: &sync.WaitGroup{}}
}

func (rd *RuneAddReader) ReadRune() (r rune, size int, err error) {
	done := make(chan struct{}, 1)
	go func() {
		rd.wg.Wait()
		select {
		case <-done:
			return
		default:
		}
		r, size, err = rd.ird.ReadRune()
		done <- struct{}{}
	}()
	select {
	case r = <-rd.ch:
		done <- struct{}{}
		size = utf8.RuneLen(r)
	case <-done:
	}
	return r, size, err
}

func (rd *RuneAddReader) Add(s []rune) {
	rd.wg.Add(len(s))
	go func() {
		for _, r := range s {
			rd.ch <- r
			rd.wg.Done()
		}
	}()
}
