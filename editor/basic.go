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
	case '(', ')', 'b':
		e.insert([]rune{')'}, to)
		e.insert([]rune{'('}, from)
	case '{', '}', 'B':
		e.insert([]rune{'}'}, to)
		e.insert([]rune{'{'}, from)
	case '[', ']', 'r':
		e.insert([]rune{']'}, to)
		e.insert([]rune{'['}, from)
	case '<', '>', 'a':
		e.insert([]rune{'>'}, to)
		e.insert([]rune{'<'}, from)
	}
}

func (e *editor) wordForward() {
	switch n := len(e.buf) - e.pos; {
	case n < 1:
		return
	case n == 1:
		e.pos = len(e.buf)
		return
	}
	switch r := e.buf[e.pos]; {
	case isWhitespace(r):
		if i := e.indexFunc(isWhitespace, e.pos+1, false); i > 0 {
			e.pos = i
			return
		}
	case isKeyword(r):
		if i := e.indexFunc(isKeyword, e.pos+1, false); i > 0 {
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
		if i := e.indexFunc(isSymbol, e.pos+1, false); i > 0 {
			if isKeyword(e.buf[i]) {
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
	switch r := e.buf[n]; {
	case isWhitespace(r):
		n = e.lastIndexFunc(isWhitespace, n, false)
		if n < 0 {
			e.pos = 0
			return
		}
	}

	switch r := e.buf[n]; {
	case isKeyword(r):
		if i := e.lastIndexFunc(isKeyword, n, false); i >= 0 {
			e.pos = i + 1
			return
		}
	default:
		if i := e.lastIndexFunc(isSymbol, n, false); i >= 0 {
			e.pos = i + 1
			return
		}
	}
	e.pos = 0
}

func (e *editor) wordForwardNonBlank() {
	i := e.indexFunc(isWhitespace, e.pos, true)
	if i < 0 {
		e.pos = len(e.buf)
		return
	}
	i = e.indexFunc(isWhitespace, i+1, false)
	if i < 0 {
		e.pos = len(e.buf)
		return
	}
	e.pos = i
}

func (e *editor) wordBackwardNonBlank() {
	i := e.lastIndexFunc(isWhitespace, e.pos, false)
	if i < 0 {
		e.pos = 0
		return
	}
	i = e.lastIndexFunc(isWhitespace, i, true)
	if i < 0 {
		e.pos = 0
		return
	}
	e.pos = i + 1
}

func (e *editor) wordEnd() {
	switch n := len(e.buf) - e.pos; {
	case n < 1:
		return
	case n == 1:
		e.pos = len(e.buf)
		return
	}
	e.pos++
	switch r := e.buf[e.pos]; {
	case isWhitespace(r):
		if i := e.indexFunc(isWhitespace, e.pos+1, false); i > 0 {
			switch r := e.buf[i]; {
			case isKeyword(r):
				if i := e.indexFunc(isKeyword, i+1, false); i > 0 {
					e.pos = i - 1
					return
				}
			default:
				if i := e.indexFunc(isSymbol, i+1, false); i > 0 {
					e.pos = i - 1
					return
				}
			}
		}
	case isKeyword(r):
		if i := e.indexFunc(isKeyword, e.pos+1, false); i > 0 {
			e.pos = i - 1
			return
		}
	default:
		if i := e.indexFunc(isSymbol, e.pos+1, false); i > 0 {
			e.pos = i - 1
			return
		}
	}
	e.pos = len(e.buf) - 1
}

func (e *editor) wordEndNonBlank() {
	switch n := len(e.buf) - e.pos; {
	case n < 1:
		return
	case n == 1:
		e.pos = len(e.buf)
		return
	}
	e.pos++
	switch r := e.buf[e.pos]; {
	case isWhitespace(r):
		if i := e.indexFunc(isWhitespace, e.pos+1, false); i > 0 {
			if i := e.indexFunc(isWhitespace, i+1, true); i > 0 {
				e.pos = i - 1
				return
			}
		}
	default:
		if i := e.indexFunc(isWhitespace, e.pos+1, true); i > 0 {
			e.pos = i - 1
			return
		}
	}
	e.pos = len(e.buf) - 1
}

func (e *editor) wordEndBackward() {
	switch n := e.pos; {
	case n < 1:
		return
	case n == 1:
		e.pos = 0
		return
	}
	switch r := e.buf[e.pos]; {
	case isWhitespace(r):
		if i := e.lastIndexFunc(isWhitespace, e.pos, false); i > 0 {
			e.pos = i
			return
		}
	case isKeyword(r):
		if i := e.lastIndexFunc(isKeyword, e.pos, false); i > 0 {
			switch {
			case isWhitespace(e.buf[i]):
				if i := e.lastIndexFunc(isWhitespace, i, false); i > 0 {
					e.pos = i
					return
				}
			default:
				e.pos = i
				return
			}
		}
	default:
		if i := e.lastIndexFunc(isSymbol, e.pos, false); i > 0 {
			switch {
			case isWhitespace(e.buf[i]):
				if i := e.lastIndexFunc(isWhitespace, i, false); i > 0 {
					e.pos = i
					return
				}
			default:
				e.pos = i
				return
			}
		}
	}
	e.pos = 0
}

func (e *editor) wordEndBackwardNonBlank() {
	switch n := e.pos; {
	case n < 1:
		return
	case n == 1:
		e.pos = 0
		return
	}
	switch r := e.buf[e.pos]; {
	case isWhitespace(r):
		if i := e.lastIndexFunc(isWhitespace, e.pos, false); i > 0 {
			e.pos = i
			return
		}
	default:
		if i := e.lastIndexFunc(isWhitespace, e.pos, true); i > 0 {
			if i := e.lastIndexFunc(isWhitespace, i, false); i > 0 {
				e.pos = i
				return
			}
		}
	}
	e.pos = 0
}
