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

	replaceMode bool
	replacedBuf []rune
}

func newInsert(ss streamSet, e *editor, s screen.Screen, conf *config.Config) *insert {
	return &insert{
		streamSet: ss,
		editor:    e,
		s:         s,
		conf:      conf,
	}
}

func newReplace(s streamSet, e *editor) *insert {
	buf := e.buf
	e.buf = nil
	return &insert{
		streamSet:   s,
		editor:      e,
		replaceMode: true,
		replacedBuf: buf,
	}
}

func (e *insert) Mode() mode {
	return modeInsert
}

func (e *insert) Run() (end continuity, next modeChanger, err error) {
	r, _, err := e.in.ReadRune()
	if err != nil {
		return end, next, err
	}
start:
	switch r {
	case CharEscape, CharCtrlC:
		if e.replaceMode {
			e.buf = e.overwrite(e.replacedBuf, e.buf, e.pos-len(e.buf))
		}
		e.move(e.pos - 1)
		next = norm()
		if e.needSave {
			e.undoTree.add(e.buf)
		}
	case CharBackspace, CharCtrlH:
		e.deleteChar()
		e.needSave = true
	case CharCtrlM, CharCtrlJ:
		end = execute
		e.needSave = true
	case CharCtrlX:
		r1, _, err := e.in.ReadRune()
		if err != nil {
			return end, next, err
		}
		r2, err := e.ctrlX(r1)
		if err != nil {
			return end, next, err
		}
		r = r2
		goto start
	case CharCtrlW:
		e.deleteWord()
		e.needSave = true
	case CharCtrlU:
		e.deleteToBeginning()
		e.needSave = true
	case CharCtrlR:
		r1, _, err := e.in.ReadRune()
		if err != nil {
			return end, next, err
		}
		e.put(r1, e.pos)
		e.needSave = true
	default:
		if r <= CharCtrlZ {
			return
		}
		e.insert([]rune{r}, e.pos)
		e.needSave = true
	}
	return
}

func (e *insert) Runes() []rune {
	if e.replaceMode {
		return e.overwrite(e.replacedBuf, e.editor.buf, e.editor.pos-len(e.buf))
	}
	return e.editor.buf
}

func (e *insert) Position() int {
	return e.editor.pos
}

func (e *insert) Message() []rune {
	if e.replaceMode {
		return []rune("-- REPLACE --")
	}
	return []rune("-- INSERT --")
}

func (e *insert) Highlight() *screen.Hi {
	return nil
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
		e.s.Refresh(e.conf, false, e.buf, e.pos, e.Highlight())
		r1, _, err := e.in.ReadRune()
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
			r2, _, err := e.in.ReadRune()
			return r2, err
		case CharCtrlE:
			e.delete(e.pos, e.pos-len(list[n1]))
			e.s.Refresh(e.conf, false, e.buf, e.pos, e.Highlight())
			r2, _, err := e.in.ReadRune()
			return r2, err
		default:
			return r1, nil
		}
		e.delete(e.pos, e.pos-len(list[n1]))
		e.insert([]rune(list[n]), e.pos)
	}
}

func (e *insert) deleteChar() {
	if !e.replaceMode {
		e.delete(e.pos-1, e.pos)
		return
	}
	if e.pos == 0 {
		return
	}
	e.pos--
	if len(e.buf) > 0 {
		e.buf = e.buf[:len(e.buf)-1]
	}
}

func (e *insert) deleteWord() {
	if !e.replaceMode {
		pos := e.pos
		e.wordBackward()
		e.delete(pos, e.pos)
		return
	}
	if e.pos == 0 {
		return
	}
	off := e.pos - len(e.buf)
	e.pos -= off
	pos := e.pos
	e.wordBackward()
	if len(e.buf) > 0 {
		e.delete(pos, e.pos)
		e.pos += off
	}
}

func (e *insert) deleteToBeginning() {
	if !e.replaceMode {
		e.delete(0, e.pos)
		return
	}
	if e.pos == 0 {
		return
	}
	off := e.pos - len(e.buf)
	e.pos -= off
	if len(e.buf) > 0 {
		e.delete(0, e.pos)
		e.pos += off
	}
}
