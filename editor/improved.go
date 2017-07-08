package editor

import (
	"fmt"
	"strings"

	"github.com/elpinal/coco3/editor/register"
)

type editor struct {
	basic
	register.Registers

	history [][]rune
	age     int
}

func (e *editor) yank(r rune, from, to int) {
	s := e.slice(from, to)
	e.Register(r, s)
}

func (e *editor) put(r rune, at int) {
	s := e.Read(r)
	e.insert(s, at)
}

func isKeyword(ch rune) bool {
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
	case isKeyword(ch):
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
		if i := e.indexFunc(func(r rune) bool { return isWhitespace(r) || isKeyword(r) }, e.pos+1, true); i > 0 {
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
	switch ch := e.buf[n]; {
	case isWhitespace(ch):
		n = e.lastIndexFunc(isWhitespace, n, false)
		if n < 0 {
			e.pos = 0
			return
		}
	}

	switch ch := e.buf[n]; {
	case isKeyword(ch):
		if i := e.lastIndexFunc(isKeyword, n, false); i >= 0 {
			e.pos = i + 1
			return
		}
	default:
		for i := n - 1; i >= 0; i-- {
			switch ch := e.buf[i]; {
			case isKeyword(ch), isWhitespace(ch):
				e.pos = i + 1
				return
			}
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
	switch ch := e.buf[e.pos]; {
	case isWhitespace(ch):
		if i := e.indexFunc(isWhitespace, e.pos+1, false); i > 0 {
			switch ch := e.buf[i]; {
			case isKeyword(ch):
				if i := e.indexFunc(isKeyword, i+1, false); i > 0 {
					e.pos = i - 1
					return
				}
			default:
				if i := e.indexFunc(func(r rune) bool { return !isWhitespace(r) && !isKeyword(r) }, i+1, false); i > 0 {
					e.pos = i - 1
					return
				}
			}
		}
	case isKeyword(ch):
		if i := e.indexFunc(isKeyword, e.pos+1, false); i > 0 {
			e.pos = i - 1
			return
		}
	default:
		if i := e.indexFunc(func(r rune) bool { return !isWhitespace(r) && !isKeyword(r) }, e.pos+1, false); i > 0 {
			e.pos = i - 1
			return
		}
	}
	e.pos = len(e.buf)
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
	switch ch := e.buf[e.pos]; {
	case isWhitespace(ch):
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
	e.pos = len(e.buf)
}

func (e *editor) toUpper(from, to int) {
	at := constrain(min(from, to), 0, len(e.buf))
	e.replace([]rune(strings.ToUpper(string(e.slice(from, to)))), at)
}

func (e *editor) toLower(from, to int) {
	at := constrain(min(from, to), 0, len(e.buf))
	e.replace([]rune(strings.ToLower(string(e.slice(from, to)))), at)
}

func (e *editor) currentWord(include bool) (from, to int) {
	if len(e.buf) == 0 {
		return 0, 0
	}
	f := func(r rune) bool { return !(isKeyword(r) || isWhitespace(r)) }
	switch ch := e.buf[e.pos]; {
	case isWhitespace(ch):
		f = isWhitespace
	case isKeyword(ch):
		f = isKeyword
	}
	from = e.lastIndexFunc(f, e.pos, false) + 1
	to = e.indexFunc(f, e.pos, false)
	if to < 0 {
		to = len(e.buf)
	}
	if include && to < len(e.buf) && isWhitespace(e.buf[to]) {
		to++
		return
	}
	if include && from > 0 && isWhitespace(e.buf[from-1]) {
		from--
		return
	}
	return
}
func (e *editor) currentWordNonBlank(include bool) (from, to int) {
	if len(e.buf) == 0 {
		return 0, 0
	}
	f := func(r rune) bool { return !isWhitespace(r) }
	if isWhitespace(e.buf[e.pos]) {
		f = isWhitespace
	}
	from = e.lastIndexFunc(f, e.pos, false) + 1
	to = e.indexFunc(f, e.pos, false)
	if to < 0 {
		to = len(e.buf)
	}
	if include && to < len(e.buf) && isWhitespace(e.buf[to]) {
		to++
		return
	}
	if include && from > 0 && isWhitespace(e.buf[from-1]) {
		from--
		return
	}
	return
}

func (e *editor) currentQuote(include bool, quote rune) (from, to int) {
	if len(e.buf) == 0 {
		return
	}
	if e.buf[e.pos] == quote {
		n := strings.Count(string(e.buf[:e.pos]), string(quote))
		if n%2 == 0 {
			// expect `to` as the position of the even-numbered quote
			to = e.index(quote, e.pos+1)
			from = e.pos
		} else {
			// expect `to` as the position of the odd-numbered quote
			from = e.lastIndex(quote, e.pos)
			to = e.pos
		}
	} else {
		from = e.lastIndex(quote, e.pos)
		if from < 0 {
			return
		}
		to = e.index(quote, e.pos)
	}
	if to < 0 {
		return
	}
	if include {
		to++
		if to < len(e.buf) && isWhitespace(e.buf[to]) {
			to++
			return
		}
		if from > 0 && isWhitespace(e.buf[from-1]) {
			from--
		}
		return
	}
	from++
	return
}

func (e *editor) charSearch(r rune) (int, error) {
	i := strings.IndexRune(string(e.slice(e.pos+1, len(e.buf))), r)
	if i < 0 {
		return 0, fmt.Errorf("pattern not found: %c", r)
	}
	return e.pos + 1 + i, nil
}

func (e *editor) charSearchBackward(r rune) (int, error) {
	i := strings.LastIndex(string(e.slice(0, e.pos)), string(r))
	if i < 0 {
		return 0, fmt.Errorf("pattern not found: %c", r)
	}
	return i, nil
}
