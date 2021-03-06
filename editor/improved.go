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

func isKeyword(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' || r == '_' || 192 <= r && r <= 255
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isSymbol(r rune) bool {
	return !(isKeyword(r) || isWhitespace(r))
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
	f := isSymbol
	switch r := e.buf[e.pos]; {
	case isWhitespace(r):
		f = isWhitespace
	case isKeyword(r):
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
	var truth bool
	if !isWhitespace(e.buf[e.pos]) {
		truth = true
	}
	from = e.lastIndexFunc(isWhitespace, e.pos, truth) + 1
	to = e.indexFunc(isWhitespace, e.pos, truth)
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
