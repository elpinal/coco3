package editor

import (
	"fmt"

	"github.com/elpinal/coco3/screen"
)

type searchType int

const (
	searchForward searchType = iota
	searchBackward
	searchHistoryForward
	searchHistoryBackward
)

type search struct {
	streamSet
	*editor

	basic *basic
	st    searchType
}

func newSearch(s streamSet, e *editor, st searchType) *search {
	return &search{
		streamSet: s,
		editor:    e,
		basic:     &basic{},
		st:        st,
	}
}

func (se *search) Mode() mode {
	return modeSearch
}

func (se *search) Position() int {
	return se.basic.pos + 1
}

func (se *search) Runes() []rune {
	return se.buf
}

func (se *search) Message() []rune {
	return append([]rune{'/'}, se.basic.buf...)
}

func (se *search) Highlight() *screen.Hi {
	return nil
}

func (se *search) Run() (end continuity, next modeChanger, err error) {
	r, _, err := se.in.ReadRune()
	if err != nil {
		return end, next, err
	}
	switch r {
	case CharCtrlM, CharCtrlJ:
	case CharEscape, CharCtrlC:
		next = norm()
		return end, next, err
	case CharBackspace, CharCtrlH:
		if len(se.basic.buf) == 0 {
			next = norm()
			return
		}
		se.basic.delete(se.basic.pos-1, se.basic.pos)
	case CharCtrlB:
		se.basic.move(0)
	case CharCtrlE:
		se.basic.move(len(se.basic.buf))
	case CharCtrlU:
		se.basic.delete(0, se.basic.pos)
	case CharCtrlW:
		// FIXME: It's redundant.
		ed := newEditor()
		ed.pos = se.basic.pos
		ed.buf = se.basic.buf
		pos := ed.pos
		ed.wordBackward()
		se.basic.delete(pos, ed.pos)
		return
	default:
		se.basic.insert([]rune{r}, se.basic.pos)
	}
	if r != CharCtrlM && r != CharCtrlJ {
		return
	}
	next = norm()
	s := string(se.basic.buf)
	if s == "" {
		return
	}
	i, err := se.search(s)
	if err != nil {
		return end, next, err
	}
	se.move(i)
	return
}

func (se *search) search(s string) (int, error) {
	found := se.editor.search(s)
	if !found {
		return 0, fmt.Errorf("pattern not found: %q", s)
	}
	if se.st == searchBackward {
		return se.previous(), nil
	}
	return se.next(), nil
}
