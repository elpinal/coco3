package main

type editor struct {
	basicEditor
	registers
}

func (e *editor) yank(r rune, from, to int) {
	s := e.slice(from, to)
	e.register(r, s)
}

func (e *editor) put(r rune, at int) {
	s := e.read(r)
	e.insert(s, at)
}

func iskeyword(ch rune) bool {
	if 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || '0' <= ch && ch <= '9' || ch == '_' || 192 <= ch && ch <= 255 {
		return true
	}
	return false
}

func isWhitespace(ch rune) bool {
	if ch == ' ' || ch == '\t' {
		return true
	}
	return false
}

func (e *editor) wordForward() {
	switch n := len(e.buf) - e.pos; {
	case n < 1:
		return
	case n == 1:
		e.pos = len(e.buf)
		return
	}
	switch ch := e.buf[e.pos]; {
	case isWhitespace(ch):
		if i := e.indexFunc(isWhitespace, e.pos+1, false); i > 0 {
			e.pos = i
			return
		}
	case iskeyword(ch):
		if i := e.indexFunc(iskeyword, e.pos+1, false); i > 0 {
			if !isWhitespace(e.buf[i]) {
				e.pos = i
				return
			}
			if i := e.indexFunc(isWhitespace, i+1, false); i > 0 {
				e.pos = i
				return
			}
		}
	default:
		if i := e.indexFunc(func(r rune) bool { return isWhitespace(r) || iskeyword(r) }, e.pos+1, true); i > 0 {
			if iskeyword(e.buf[i]) {
				e.pos = i
				return
			}
			if i := e.indexFunc(isWhitespace, i+1, false); i > 0 {
				e.pos = i
				return
			}
		}
	}
	e.pos = len(e.buf)
}

func (e *editor) wordBackward() {
	switch e.pos {
	case 0:
		return
	case 1:
		e.pos = 0
		return
	}

	n := e.pos - 1
	switch ch := e.buf[n]; {
	case isWhitespace(ch):
		n = e.lastIndexFunc(isWhitespace, n, false)
		if n < 0 {
			e.pos = 0
			return
		}
	}

	switch ch := e.buf[n]; {
	case iskeyword(ch):
		if i := e.lastIndexFunc(iskeyword, n, false); i >= 0 {
			e.pos = i + 1
			return
		}
	default:
		for i := n - 1; i >= 0; i-- {
			switch ch := e.buf[i]; {
			case iskeyword(ch), isWhitespace(ch):
				e.pos = i + 1
				return
			}
		}
	}
	e.pos = 0
}
