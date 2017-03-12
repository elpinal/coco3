package editor

type insert struct {
	streamSet
	*editor
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
	switch r {
	case CharEscape:
		e.move(e.pos - 1)
		next = modeNormal
	case CharBackspace:
		e.delete(e.pos-1, e.pos)
	case CharCtrlM:
		end = true
	default:
		e.insert([]rune{r}, e.pos)
	}
	return end, next, err
}

func (e *insert) Runes() []rune {
	return e.editor.buf
}

func (e *insert) Position() int {
	return e.editor.pos
}
