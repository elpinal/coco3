package main

type basicEditor struct {
	pos int
	buf []rune
}

// move moves the position to the given position.
// Given a invalid position, move sets the position at the end of the buffer.
// Valid positions are in range [0, len(e.buf)].
func (e *basicEditor) move(to int) {
	switch {
	case to >= len(e.buf):
		e.pos = len(e.buf)
	case to <= 0:
		e.pos = 0
	default:
		e.pos = to
	}
}

// insert inserts s into the buffer at the given position.
// Given a invalid position, insert considers the position to be at the end of the buffer.
func (e *basicEditor) insert(s []rune, at int) {
	switch {
	case at < 0:
		at = 0
	case at > len(e.buf):
		at = len(e.buf)
	}
	switch at {
	case 0:
		e.buf = append(s, e.buf...)
	case len(e.buf):
		e.buf = append(e.buf, s...)
	default:
		s = append(e.buf[:at], s...)
		e.buf = append(s, e.buf[at:]...)
	}
	if at <= e.pos {
		e.pos += len(s)
	}
}

// delete deletes runes from the buffer [from, to].
// Given a invalid position, delete considers the position to be at the end of the buffer.
func (e *basicEditor) delete(from, to int) {
	left := from
	right := to
	if from > to {
		left = to
		right = from
	}
	if left < 0 {
		left = 0
	}
	if right > len(e.buf) {
		right = len(e.buf)
	}
	switch {
	case left == 0:
		e.buf = e.buf[right:]
	case right == len(e.buf):
		e.buf = e.buf[:left]
	default:
		e.buf = append(e.buf[:left], e.buf[right:]...)
	}
	switch {
	case e.pos < left:
	case right < e.pos:
		e.pos = e.pos - (right - left)
	default:
		e.pos = left
	}
}
