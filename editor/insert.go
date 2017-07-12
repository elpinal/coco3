package editor

import (
	"github.com/elpinal/coco3/complete"
	"github.com/elpinal/coco3/config"
	"github.com/elpinal/coco3/screen"
)

type insert struct {
	streamSet
	*editor
	s    screen.Screen
	conf *config.Config

	needSave bool
}

func (e *insert) Mode() mode {
	return modeInsert
}

func (e *insert) Run() (end bool, next mode, err error) {
	next = modeInsert
	r, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return end, next, err
	}
start:
	switch r {
	case CharEscape:
		e.move(e.pos - 1)
		next = modeNormal
		if e.needSave {
			e.undoTree.add(e.buf)
		}
	case CharBackspace:
		e.delete(e.pos-1, e.pos)
		e.needSave = true
	case CharCtrlM:
		end = true
		e.needSave = true
	case CharCtrlX:
		r1, _, err := e.streamSet.in.ReadRune()
		if err != nil {
			return end, next, err
		}
		r2, err := e.ctrlX(r1)
		if err != nil {
			return end, next, err
		}
		r = r2
		goto start
	default:
		e.insert([]rune{r}, e.pos)
		e.needSave = true
	}
	return end, next, err
}

func (e *insert) Runes() []rune {
	return e.editor.buf
}

func (e *insert) Position() int {
	return e.editor.pos
}

func (e *insert) ctrlX(r rune) (rune, error) {
	var f func([]rune, int) ([]string, error)
	switch r {
	case CharCtrlF:
		f = complete.File
	default:
		return r, nil
	}

	list, err := f(e.buf, e.pos)
	if err != nil {
		return 0, err
	}
	list = append(list, "")
	e.insert([]rune(list[0]), e.pos)
	e.needSave = true
	n := 0
	for {
		e.s.Refresh(e.conf, e.buf, e.pos)
		r1, _, err := e.streamSet.in.ReadRune()
		if err != nil {
			return 0, err
		}
		n1 := n
		switch r1 {
		case CharCtrlN, r:
			n++
			if len(list) <= n {
				n = 0
			}
		case CharCtrlP:
			n--
			if n < 0 {
				n = len(list) - 1
			}
		case CharCtrlY:
			r2, _, err := e.streamSet.in.ReadRune()
			return r2, err
		case CharCtrlE:
			e.delete(e.pos, e.pos-len(list[n1]))
			e.s.Refresh(e.conf, e.buf, e.pos)
			r2, _, err := e.streamSet.in.ReadRune()
			return r2, err
		default:
			return r1, nil
		}
		e.delete(e.pos, e.pos-len(list[n1]))
		e.insert([]rune(list[n]), e.pos)
	}
}
