package editor

import (
	"strings"
)

// basic represents a basic editor.
// Valid positions are in range [0, len(e.buf)].
type basic struct {
	pos int
	buf []rune
}

// move moves the position to the given position.
// Given a invalid position, move sets the position at the end of the buffer.
func (e *basic) move(to int) {
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
func (e *basic) insert(s []rune, at int) {
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
		x := append(e.buf[:at:at], s...)
		e.buf = append(x, e.buf[at:]...)
	}
	if at <= e.pos {
		e.pos += len(s)
	}
}

func max(m, n int) int {
	if m > n {
		return m
	}
	return n
}

func min(m, n int) int {
	if m < n {
		return m
	}
	return n
}

func constrain(n, low, high int) int {
	if n < low {
		return low
	}
	if n > high {
		return high
	}
	return n
}

// delete deletes runes from the buffer [from, to].
// Given a invalid position, delete considers the position to be at the end of the buffer.
func (e *basic) delete(from, to int) {
	left := constrain(min(from, to), 0, len(e.buf))
	right := constrain(max(from, to), 0, len(e.buf))
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

// slice slices the buffer [from, to].
// Given a invalid position, slice considers the position to be at the end of the buffer.
func (e *basic) slice(from, to int) []rune {
	left := constrain(min(from, to), 0, len(e.buf))
	right := constrain(max(from, to), 0, len(e.buf))
	s := make([]rune, right-left)
	copy(s, e.buf[left:right])
	return s
}

func (e *basic) index(ch rune, start int) int {
	start = constrain(start, 0, len(e.buf))
	for i := start; i < len(e.buf); i++ {
		if e.buf[i] == ch {
			return i
		}
	}
	return -1
}

func (e *basic) lastIndex(ch rune, last int) int {
	last = constrain(last, 0, len(e.buf))
	for i := last - 1; i >= 0; i-- {
		if e.buf[i] == ch {
			return i
		}
	}
	return -1
}

// indexFunc(f, start, true) == indexFunc(func(r) bool { return !f(r) }, start, false)
func (e *basic) indexFunc(f func(rune) bool, start int, truth bool) int {
	start = constrain(start, 0, len(e.buf))
	for i := start; i < len(e.buf); i++ {
		if f(e.buf[i]) == truth {
			return i
		}
	}
	return -1
}

// lastIndexFunc(f, start, true) == lastIndexFunc(func(r) bool { return !f(r) }, start, false)
func (e *basic) lastIndexFunc(f func(rune) bool, last int, truth bool) int {
	last = constrain(last, 0, len(e.buf))
	for i := last - 1; i >= 0; i-- {
		if f(e.buf[i]) == truth {
			return i
		}
	}
	return -1
}

func (e *basic) replace(s []rune, at int) {
	if s == nil {
		return
	}
	// Q. What should we do when `at` < 0 or len(e.buf) < `at`?
	switch {
	case len(e.buf) <= at:
		e.buf = append(e.buf, []rune(strings.Repeat(" ", at-len(e.buf)))...)
		e.buf = append(e.buf, s...)
	case at+len(s) <= 0:
		// no-op
	case at < 0:
		for i := 0; i < at+len(s); i++ {
			e.buf[i] = s[i-at]
		}
	case len(e.buf) < at+len(s):
		v := make([]rune, at+len(s))
		copy(v, e.buf)
		for i := at; i < at+len(s); i++ {
			v[i] = s[i-at]
		}
		e.buf = v
	default:
		for i := at; i < at+len(s); i++ {
			e.buf[i] = s[i-at]
		}
	}
}

func (e *basic) siege(from int, to int, r rune) {
	switch r {
	case '\'', '"', '`', '@', '*', '+', '_', '|', '$':
		e.insert([]rune{r}, to)
		e.insert([]rune{r}, from)
	case '(', ')':
		e.insert([]rune{')'}, to)
		e.insert([]rune{'('}, from)
	case '{', '}':
		e.insert([]rune{'}'}, to)
		e.insert([]rune{'{'}, from)
	case '[', ']':
		e.insert([]rune{']'}, to)
		e.insert([]rune{'['}, from)
	case '<', '>':
		e.insert([]rune{'>'}, to)
		e.insert([]rune{'<'}, from)
	}
}
