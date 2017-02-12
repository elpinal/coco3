package main

func (cl *commandline) toLeft() {
	if cl.index == 0 {
		return
	}
	cl.index--
}

func (cl *commandline) toRight() {
	if cl.index == len(cl.buf)-1 {
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
