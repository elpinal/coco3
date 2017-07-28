package editor

import (
	"fmt"

	"github.com/elpinal/coco3/editor/register"
	"github.com/elpinal/coco3/screen"
)

type opArg struct {
	opType     int // current operator type
	opStart    int
	opCount    int
	motionType int
	inclusive  bool
}

type normalSet struct {
	opArg
	finishOp bool
	regName  rune
}

type nvCommon struct {
	streamSet
	*editor

	count int
}

type normal struct {
	nvCommon

	normalSet
}

func newNormal(s streamSet, e *editor) *normal {
	return &normal{
		nvCommon: nvCommon{
			streamSet: s,
			editor:    e,
		},
	}
}

func norm() modeChanger {
	return func(b *balancer) (moder, error) {
		return newNormal(b.streamSet, b.editor), nil
	}
}

func (e *normal) Mode() mode {
	return modeNormal
}

func (e *normal) Run() (end continuity, next modeChanger, err error) {
	e.finishOp = e.opType != OpNop
	r, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return end, next, err
	}
	for ('1' <= r && r <= '9') || (e.count != 0 && r == '0') {
		e.count = e.count*10 + int(r-'0')
		r1, _, err := e.streamSet.in.ReadRune()
		if err != nil {
			return end, next, err
		}
		r = r1
	}
	if e.count == 0 {
		e.count = 1
	}
	if e.opCount > 0 {
		e.count *= e.opCount
	}
	if e.regName == 0 {
		e.regName = register.Unnamed
	}
	if cmd, ok := normalCommands[r]; ok {
		if m := cmd(e, r); m != nil {
			next = m
		}
		if n := e.doPendingOperator(); n != nil {
			next = n
		}
		e.count = 0
		if e.pos == len(e.buf) {
			e.move(e.pos - 1)
		}
	}
	return
}

func (e *normal) Runes() []rune {
	return e.editor.buf
}

func (e *normal) Position() int {
	return e.editor.pos
}

func (e *normal) Message() []rune {
	return nil
}

func (e *normal) Highlight() *screen.Hi {
	return nil
}

type normalCommand = func(*normal, rune) modeChanger

var normalCommands = map[rune]normalCommand{
	CharCtrlR: (*normal).redoCmd,
	'"':       (*normal).handleRegister,
	':':       (*normal).commandline,
	'|':       (*normal).column,
	'/':       (*normal).search,
	'?':       (*normal).searchBackward,
	'+':       (*normal).increment,
	'-':       (*normal).decrement,
	'$':       (*normal).endline,
	'^':       (*normal).beginlineNonBlank,
	'0':       (*normal).beginline,
	'A':       (*normal).edit,
	'B':       (*normal).wordBack,
	'C':       (*normal).abbrev,
	'D':       (*normal).abbrev,
	'E':       (*normal).word,
	'F':       (*normal).searchCharacterBackward,
	'I':       (*normal).edit,
	'N':       (*normal).previous,
	'R':       (*normal).replaceMode,
	'W':       (*normal).word,
	'X':       (*normal).abbrev,
	'Y':       (*normal).abbrev,
	'a':       (*normal).edit,
	'b':       (*normal).wordBack,
	'c':       (*normal).operator1,
	'd':       (*normal).operator1,
	'e':       (*normal).word,
	'f':       (*normal).searchCharacter,
	'g':       (*normal).gCmd,
	'h':       (*normal).left,
	'i':       (*normal).edit,
	'j':       (*normal).down,
	'k':       (*normal).up,
	'l':       (*normal).right,
	'n':       (*normal).next,
	'p':       (*normal).put1,
	'r':       (*normal).replace,
	'u':       (*normal).undoCmd,
	'v':       (*normal).visual,
	'w':       (*normal).word,
	'x':       (*normal).abbrev,
	'y':       (*normal).operator1,
}

func (e *nvCommon) endline(r rune) (_ modeChanger) {
	e.move(len(e.buf))
	return
}

func (e *nvCommon) beginline(r rune) (_ modeChanger) {
	e.move(0)
	return
}

func (e *nvCommon) beginlineNonBlank(r rune) (_ modeChanger) {
	i := e.indexFunc(isWhitespace, 0, false)
	if i < 0 {
		return e.endline(r)
	}
	e.move(i)
	return
}

func (e *nvCommon) wordBack(r rune) (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		switch r {
		case 'b':
			e.wordBackward()
		case 'B':
			e.wordBackwardNonBlank()
		}
	}
	return
}

func (e *normal) operator1(r rune) modeChanger {
	return e.operator(string(r))
}

func (e *normal) operator(s string) (_ modeChanger) {
	op := opChars[s]
	if op == e.opType { // double operator
		e.motionType = mline
	} else {
		e.opStart = e.pos
		e.opType = op
		e.opCount = e.count
	}
	return
}

func (e *nvCommon) left(r rune) (_ modeChanger) {
	e.move(e.pos - e.count)
	return
}

func ins(rightmost bool) modeChanger {
	return func(b *balancer) (moder, error) {
		if rightmost {
			// Revert to the rightmost position.
			b.pos = len(b.buf)
		}
		return newInsert(b.streamSet, b.editor, b.s, b.conf), nil
	}
}

func (e *normal) edit(r rune) modeChanger {
	if (r == 'a' || r == 'i') && e.opType != OpNop {
		e.object(r)
		return nil
	}
	switch r {
	case 'A':
		e.move(len(e.buf))
	case 'I':
		_ = e.beginlineNonBlank(r)
	case 'a':
		e.move(e.pos + 1)
	}
	return ins(e.pos == len(e.buf))
}

func (e *normal) object(r rune) {
	var include bool
	if r == 'a' {
		include = true
	}
	var from, to int
	r1, _, _ := e.streamSet.in.ReadRune()
	switch r1 {
	case 'w':
		from, to = e.currentWord(include)
	case 'W':
		from, to = e.currentWordNonBlank(include)
	case '"', '\'', '`':
		from, to = e.currentQuote(include, r1)
		if from < 0 || to < 0 {
			return
		}
	default:
		return
	}
	e.opStart = from
	e.pos = to
}

func (e *normal) down(r rune) (_ modeChanger) {
	if e.age >= len(e.history)-1 {
		return
	}
	e.history[e.age] = e.buf
	e.age++
	e.buf = e.history[e.age]
	e.pos = len(e.buf)
	return
}

func (e *normal) up(r rune) (_ modeChanger) {
	if e.age <= 0 {
		return
	}
	if e.age == len(e.history) {
		e.history = append(e.history, e.buf)
	} else {
		e.history[e.age] = e.buf
	}
	e.age--
	e.buf = e.history[e.age]
	e.pos = len(e.buf)
	return
}

func (e *nvCommon) right(r rune) (_ modeChanger) {
	e.move(e.pos + e.count)
	return
}

func (e *normal) put1(r rune) (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.put(e.regName, e.pos+1)
	}
	e.undoTree.add(e.buf)
	return
}

func (e *normal) replace(r rune) (_ modeChanger) {
	r1, _, _ := e.streamSet.in.ReadRune()
	s := make([]rune, e.count)
	for i := 0; i < e.count; i++ {
		s[i] = r1
	}
	e.editor.replace(s, e.pos)
	e.move(e.pos + e.count - 1)
	e.undoTree.add(e.buf)
	return
}

func (e *normal) word(r rune) (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		switch r {
		case 'w':
			e.wordForward()
		case 'W':
			e.wordForwardNonBlank()
		case 'e':
			e.inclusive = true
			e.wordEnd()
		case 'E':
			e.inclusive = true
			e.wordEndNonBlank()
		}
	}
	return
}

func (e *normal) doPendingOperator() (_ modeChanger) {
	if !e.finishOp {
		return
	}
	from := e.opStart
	to := e.pos
	if e.inclusive {
		to++
	}
	if e.motionType == mline {
		from = 0
		to = len(e.buf)
	}
	switch e.opType {
	case OpDelete:
		e.yank(e.regName, from, to)
		e.delete(from, to)
		e.undoTree.add(e.buf)
	case OpYank:
		e.yank(e.regName, from, to)
	case OpChange:
		e.yank(e.regName, from, to)
		e.delete(from, to)
		e.undoTree.add(e.buf)
		return ins(e.pos == len(e.buf))
	case OpLower:
		e.toLower(from, to)
		e.undoTree.add(e.buf)
	case OpUpper:
		e.toUpper(from, to)
		e.undoTree.add(e.buf)
	case OpTilde:
		e.swapCase(from, to)
		e.undoTree.add(e.buf)
	}
	e.clearOp()
	e.move(min(from, to))
	return
}

func (e *normal) clearOp() {
	e.opArg = opArg{}
}

func (e *normal) abbrev(r rune) (_ modeChanger) {
	amap := map[rune][]rune{
		'x': []rune("dl"),
		'X': []rune("dh"),
		'D': []rune("d$"),
		'C': []rune("c$"),
		'Y': []rune("y$"),
	}
	e.streamSet.in.Add(amap[r])
	e.opCount = e.count
	return
}

func (e *nvCommon) searchCharacter(r rune) (_ modeChanger) {
	r1, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return
	}
	pos := e.pos
	for i := 0; i < e.count; i++ {
		i, err := e.charSearch(r1)
		if err != nil {
			e.move(pos)
			return
		}
		e.move(i)
	}
	return
}

func (e *nvCommon) searchCharacterBackward(r rune) (_ modeChanger) {
	r1, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return
	}
	pos := e.pos
	for i := 0; i < e.count; i++ {
		i, err := e.charSearchBackward(r1)
		if err != nil {
			e.move(pos)
			return
		}
		e.move(i)
	}
	return
}

func (e *normal) gCmd(r rune) (_ modeChanger) {
	r1, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return
	}
	switch r1 {
	case 'u', 'U', '~':
		return e.operator(string([]rune{r, r1}))
	case '/':
		return e.searchHistory(r)
	}
	return
}

func (e *normal) undoCmd(r rune) (_ modeChanger) {
	e.undo()
	return
}

func (e *normal) redoCmd(r rune) (_ modeChanger) {
	e.redo()
	return
}

func (e *normal) replaceMode(r rune) modeChanger {
	return func(b *balancer) (moder, error) {
		return newReplace(b.streamSet, b.editor), nil
	}
}

func (e *normal) handleRegister(r rune) (_ modeChanger) {
	r1, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return
	}
	if !register.IsValid(r1) {
		return
	}
	e.regName = r1
	return
}

func (e *normal) commandline(r rune) modeChanger {
	return func(b *balancer) (moder, error) {
		return newCommandline(b.streamSet, b.editor), nil
	}
}

func (e *normal) visual(r rune) modeChanger {
	return func(b *balancer) (moder, error) {
		return newVisual(b.streamSet, b.editor), nil
	}
}

func (e *nvCommon) column(_ rune) (_ modeChanger) {
	e.move(constrain(e.count-1, 0, len(e.buf)))
	return
}

func (e *normal) search(_ rune) modeChanger {
	return func(b *balancer) (moder, error) {
		return newSearch(b.streamSet, b.editor, searchForward), nil
	}
}

func (e *normal) searchBackward(_ rune) modeChanger {
	return func(b *balancer) (moder, error) {
		return newSearch(b.streamSet, b.editor, searchBackward), nil
	}
}

func (e *normal) next(_ rune) (_ modeChanger) {
	e.move(e.nvCommon.next())
	return
}

func (e *normal) previous(_ rune) (_ modeChanger) {
	e.move(e.nvCommon.previous())
	return
}

func (e *normal) searchHistory(_ rune) (_ modeChanger) {
	return func(b *balancer) (moder, error) {
		return newSearch(b.streamSet, b.editor, searchHistoryForward), nil
	}
}

func (e *normal) indexNumber() int {
	for i, r := range e.buf[e.pos:] {
		if '0' <= r && r <= '9' {
			return i + e.pos
		}
	}
	return -1
}

func (e *normal) parseNumber(i int) (a int, l int) {
	for n := i; n < len(e.buf); n++ {
		r := int(e.buf[n])
		if '0' <= r && r <= '9' {
			a = 10*a + r - '0'
			continue
		}
		return a, n - i
	}
	return a, len(e.buf) - i
}

func (e *normal) increment(_ rune) (_ modeChanger) {
	i := e.indexNumber()
	if i < 0 {
		return
	}
	n, l := e.parseNumber(i)
	e.delete(i, i+l)
	e.insert([]rune(fmt.Sprint(n+1)), i)
	e.move(i)
	return
}

func (e *normal) decrement(_ rune) (_ modeChanger) {
	i := e.indexNumber()
	if i < 0 {
		return
	}
	n, l := e.parseNumber(i)
	e.delete(i, i+l)
	e.insert([]rune(fmt.Sprint(n-1)), i)
	e.move(i)
	return
}
