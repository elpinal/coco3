package editor

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/elpinal/coco3/editor/register"
	"github.com/elpinal/revim"
)

type searchRange [][]int

type editor struct {
	basic
	register.Registers
	undoTree

	history [][]rune
	age     int

	sp string // search pattern
	sr searchRange
}

func newEditor() *editor {
	r := register.Registers{}
	r.Init()
	return &editor{
		undoTree:  newUndoTree(),
		Registers: r,
	}
}

func newEditorBuffer(buf []rune) *editor {
	e := newEditor()
	e.buf = buf
	return e
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
	switch ch := e.buf[e.pos]; {
	case isWhitespace(ch):
		if i := e.lastIndexFunc(isWhitespace, e.pos, false); i > 0 {
			e.pos = i
			return
		}
	case isKeyword(ch):
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
		if i := e.lastIndexFunc(func(r rune) bool { return isWhitespace(r) || isKeyword(r) }, e.pos, true); i > 0 {
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
	switch ch := e.buf[e.pos]; {
	case isWhitespace(ch):
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

func (e *editor) toUpper(from, to int) {
	at := constrain(min(from, to), 0, len(e.buf))
	for i := at; i < constrain(max(from, to), 0, len(e.buf)); i++ {
		e.buf[i] = unicode.ToUpper(e.buf[i])
	}
}

func (e *editor) toLower(from, to int) {
	at := constrain(min(from, to), 0, len(e.buf))
	for i := at; i < constrain(max(from, to), 0, len(e.buf)); i++ {
		e.buf[i] = unicode.ToLower(e.buf[i])
	}
}

func (e *editor) switchCase(from, to int) {
	at := constrain(min(from, to), 0, len(e.buf))
	for i := at; i < constrain(max(from, to), 0, len(e.buf)); i++ {
		if unicode.IsLower(e.buf[i]) {
			e.buf[i] = unicode.ToUpper(e.buf[i])
		} else if unicode.IsUpper(e.buf[i]) {
			e.buf[i] = unicode.ToLower(e.buf[i])
		}
	}
}

func (e *editor) currentWord(include bool) (from, to int) {
	if len(e.buf) == 0 {
		return 0, 0
	}
	if e.pos == len(e.buf) {
		return -1, -1
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
	if e.pos == len(e.buf) {
		return -1, -1
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
	if e.pos == len(e.buf) {
		return -1, -1
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
			to = -1
			return -1, -1
		}
		to = e.index(quote, e.pos)
	}
	if to < 0 {
		return -1, -1
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

func (e *editor) currentParen(include bool, p1, p2 rune) (from, to int) {
	if len(e.buf) == 0 {
		return
	}
	if e.pos == len(e.buf) {
		return -1, -1
	}

	from = e.searchLeft(p1, p2)
	if from < 0 {
		return -1, -1
	}
	to = e.searchRight(p1, p2)
	if to < 0 {
		return -1, -1
	}
	if include {
		to++
		return
	}
	from++
	return
}

func (e *editor) searchLeft(p1, p2 rune) int {
	if e.pos != len(e.buf) && e.buf[e.pos] == p1 {
		return e.pos
	}
	var level int
	for i := e.pos - 1; i >= 0; i-- {
		switch e.buf[i] {
		case p1:
			if level == 0 {
				return i
			}
			level--
		case p2:
			level++
		}
	}
	return -1
}

func (e *editor) searchRight(p1, p2 rune) int {
	if e.pos != len(e.buf) && e.buf[e.pos] == p2 {
		return e.pos
	}
	var level int
	for i := e.pos + 1; i < len(e.buf); i++ {
		switch e.buf[i] {
		case p2:
			if level == 0 {
				return i
			}
			level--
		case p1:
			level++
		}
	}
	return -1
}

func (e *editor) charSearch(r rune) (int, error) {
	for i := e.pos + 1; i < len(e.buf); i++ {
		if e.buf[i] == r {
			return i, nil
		}
	}
	return 0, fmt.Errorf("pattern not found: %c", r)
}

func (e *editor) charSearchBefore(r rune) (int, error) {
	for i := e.pos + 1; i < len(e.buf); i++ {
		if e.buf[i] == r {
			return i - 1, nil
		}
	}
	return 0, fmt.Errorf("pattern not found: %c", r)
}

func (e *editor) charSearchBackward(r rune) (int, error) {
	for i := e.pos - 1; i >= 0; i-- {
		if e.buf[i] == r {
			return i, nil
		}
	}
	return 0, fmt.Errorf("pattern not found: %c", r)
}

func (e *editor) charSearchBackwardAfter(r rune) (int, error) {
	for i := e.pos - 1; i >= 0; i-- {
		if e.buf[i] == r {
			return i + 1, nil
		}
	}
	return 0, fmt.Errorf("pattern not found: %c", r)
}

func (e *editor) undo() {
	s, ok := e.undoTree.undo()
	if !ok {
		return
	}
	e.buf = make([]rune, len(s))
	copy(e.buf, s)
	e.move(0)
}

func (e *editor) redo() {
	s, ok := e.undoTree.redo()
	if !ok {
		return
	}
	e.buf = make([]rune, len(s))
	copy(e.buf, s)
	e.move(0)
}

func (e *editor) overwrite(base []rune, cover []rune, at int) []rune {
	n := constrain(at, 0, len(base))
	s := make([]rune, max(len(base), n+len(cover)))
	copy(s[:n], base)
	copy(s[n:], cover)
	if n+len(cover) < len(base) {
		copy(s[n+len(cover):], base[n+len(cover):])
	}
	return s
}

func (e *editor) search(s string) (found bool) {
	e.sp = s
	e.sr = e.sr[:0]
	if s == "" {
		return false
	}
	re, err := revim.Compile(s)
	if err != nil {
		// TODO: report error
		return false
	}
	loc := re.FindAllStringIndex(string(e.buf))
	if loc == nil {
		return false
	}
	e.sr = loc
	return true
}

func (e *editor) next() int {
	found := e.search(e.sp)
	if !found {
		return e.pos
	}
	for _, sr := range e.sr {
		i := sr[0]
		if i > e.pos {
			return i
		}
	}
	return e.sr[0][0]
}

func (e *editor) previous() int {
	found := e.search(e.sp)
	if !found {
		return e.pos
	}
	for n := len(e.sr) - 1; 0 <= n; n-- {
		i := e.sr[n][0]
		if e.pos > i {
			return i
		}
	}
	return e.sr[len(e.sr)-1][0]
}
