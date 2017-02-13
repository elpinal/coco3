package main

type basicEditor struct {
	pos int
	buf []rune
}

// move moves the position.
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
