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
			next = cmd.fn(e, r)
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
	r   rune                     // first char
	fn  func(*normal, rune) mode // function for this command
	arg int
}

var normalCommands = []normalCommand{
	{'$', (*normal).endline, 0},
	{'0', (*normal).beginline, 0},
	{'A', (*normal).edit, 0},
	{'B', (*normal).wordBack, 0},
	{'I', (*normal).edit, 0},
	{'W', (*normal).word, 0},
	{'a', (*normal).edit, 0},
	{'b', (*normal).wordBack, 0},
	{'h', (*normal).left, 0},
	{'i', (*normal).edit, 0},
	{'l', (*normal).right, 0},
	{'w', (*normal).word, 0},
}

func (e *normal) endline(r rune) mode {
	e.move(len(e.buf) - 1)
	return modeNormal
}

func (e *normal) beginline(r rune) mode {
	e.move(0)
	return modeNormal
}

func (e *normal) wordBack(r rune) mode {
	switch r {
	case 'b':
		e.wordBackward()
	case 'B':
		e.wordBackwardNonBlank()
	}
	return modeNormal
}

func (e *normal) left(r rune) mode {
	e.move(e.pos - 1)
	return modeNormal
}

func (e *normal) edit(r rune) mode {
	switch r {
	case 'A':
		e.move(len(e.buf))
	case 'I':
		e.move(0)
	case 'a':
		e.move(e.pos + 1)
	}
	return modeInsert
}

func (e *normal) right(r rune) mode {
	e.move(e.pos + 1)
	return modeNormal
}

func (e *normal) word(r rune) mode {
	switch r {
	case 'w':
		e.wordForward()
	case 'W':
		e.wordForwardNonBlank()
	}
	return modeNormal
}