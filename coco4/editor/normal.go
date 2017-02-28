package editor

type normal struct {
	streamSet
	*editor
}

func (e *normal) Mode() mode {
	return modeNormal
}

func (e *normal) Run() (end bool, next mode, err error) {
	next = modeNormal
	r, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return end, next, err
	}
	for _, cmd := range normalCommands {
		if cmd.r == r {
			next = cmd.fn(e)
			return
		}
	}
	return end, next, err
}

func (e *normal) Runes() []rune {
	return e.editor.buf
}

func (e *normal) Position() int {
	return e.editor.pos
}

type normalCommand struct {
	r   rune               // first char
	fn  func(*normal) mode // function for this command
	arg int
}

var normalCommands = []normalCommand{
	{'$', (*normal).endline, 0},
	{'0', (*normal).beginline, 0},
	{'b', (*normal).wordBack, 0},
	{'h', (*normal).left, 0},
	{'i', (*normal).edit, 0},
	{'l', (*normal).right, 0},
	{'w', (*normal).word, 0},
}

func (e *normal) endline() mode {
	e.move(len(e.buf) - 1)
	return modeNormal
}

func (e *normal) beginline() mode {
	e.move(0)
	return modeNormal
}

func (e *normal) wordBack() mode {
	e.wordBackward()
	return modeNormal
}

func (e *normal) left() mode {
	e.move(e.pos - 1)
	return modeNormal
}

func (e *normal) edit() mode {
	return modeInsert
}

func (e *normal) right() mode {
	e.move(e.pos + 1)
	return modeNormal
}

func (e *normal) word() mode {
	e.wordForward()
	return modeNormal
}
