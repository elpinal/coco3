package main

func (cl *commandline) toLeft() {
	if cl.index == 0 {
		return
	}
	cl.index--
}

func (cl *commandline) toRight() {
	if cl.index >= len(cl.buf)-1 {
		return
	}
	cl.index++
}

func (cl *commandline) toTheFirst() {
	cl.index = 0
}

func (cl *commandline) toTheFirstNonBlank() {
	for i, ch := range cl.buf {
		if ch != ' ' && ch != '\t' {
			cl.index = i
		}
	}
	cl.index = len(cl.buf) - 1
}

func (cl *commandline) toTheEnd() {
	cl.index = len(cl.buf) - 1
}

func (cl *commandline) prevHistory() {
	if cl.hist.i == 0 {
		return
	}
	if cl.hist.i == len(cl.hist.lines) {
		cl.hist.lines = append(cl.hist.lines, cl.buf)
	} else {
		cl.hist.lines[cl.hist.i] = cl.buf
	}
	cl.hist.i--
	cl.buf = cl.hist.lines[cl.hist.i]
	if len(cl.buf) == 0 {
		cl.index = 0
		return
	}
	cl.index = len(cl.buf) - 1
}

func (cl *commandline) nextHistory() {
	if cl.hist.i >= len(cl.hist.lines)-1 {
		return
	}
	cl.hist.lines[cl.hist.i] = cl.buf
	cl.hist.i++
	cl.buf = cl.hist.lines[cl.hist.i]
	if len(cl.buf) == 0 {
		cl.index = 0
		return
	}
	cl.index = len(cl.buf) - 1
}

func iskeyword(ch rune) bool {
	if 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || '0' <= ch && ch <= '9' ||
		ch == '_' || 192 <= ch && ch <= 255 {
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

func (cl *commandline) wordForward() {
	if len(cl.buf[cl.index:]) <= 1 {
		return
	}
	n := -1
	switch ch1 := cl.buf[cl.index]; {
	case isWhitespace(ch1):
		for i, ch := range cl.buf[cl.index+1:] {
			if !isWhitespace(ch) {
				cl.index += i + 1
				return
			}
		}
	case iskeyword(ch1):
		for i, ch := range cl.buf[cl.index+1:] {
			if !iskeyword(ch) {
				if !isWhitespace(ch) {
					cl.index += i + 1
					return
				}
				n = i + 1
				break
			}
		}
	default:
		for i, ch := range cl.buf[cl.index+1:] {
			switch {
			case iskeyword(ch):
				cl.index += i + 1
				return
			case isWhitespace(ch):
				n = i + 1
			}
		}
	}

	if n == -1 {
		cl.index = len(cl.buf) - 1
		return
	}

	for i, ch := range cl.buf[cl.index+n:] {
		if !isWhitespace(ch) {
			cl.index = cl.index + n + i
			return
		}
	}

	cl.index = len(cl.buf) - 1
}

func (cl *commandline) deleteUnder() {
	switch cl.index {
	case 0:
		cl.buf = cl.buf[1:]
	case len(cl.buf) - 1:
		cl.buf = cl.buf[:cl.index]
		cl.index--
	default:
		cl.buf = append(cl.buf[:cl.index], cl.buf[cl.index+1:]...)
	}
}
