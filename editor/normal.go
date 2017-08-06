package editor

import (
	"fmt"

	"github.com/elpinal/coco3/editor/register"
	"github.com/elpinal/coco3/screen"
)

type nvCommon struct {
	streamSet
	*editor

	count int
}

type normal struct {
	nvCommon

	regName rune
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
	if e.regName == 0 {
		e.regName = register.Unnamed
	}
	if cmd, ok := normalCommands[r]; ok {
		if m := cmd(e); m != nil {
			next = m
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

type normalCommand = func(*normal) modeChanger

var normalCommands = map[rune]normalCommand{
	CharCtrlR: (*normal).redoCmd,
	'"':       (*normal).handleRegister,
	':':       (*normal).commandline,
	'|':       (*normal).column,
	'/':       (*normal).search,
	'?':       (*normal).searchBackward,
	'+':       (*normal).increment,
	'-':       (*normal).decrement,
	'~':       (*normal).switchCase,
	'[':       (*normal).prevUnmatched,
	']':       (*normal).nextUnmatched,
	'%':       (*normal).moveToMatch,
	'$':       (*normal).endline,
	'^':       (*normal).beginlineNonBlank,
	'0':       (*normal).beginline,
	'A':       (*normal).appendAtEnd,
	'B':       (*normal).wordBackNonBlank,
	'C':       (*normal).changeToEnd,
	'D':       (*normal).deleteToEnd,
	'E':       (*normal).wordEndNonBlank,
	'F':       (*normal).searchCharacterBackward,
	'I':       (*normal).insertFirstNonBlank,
	'N':       (*normal).previous,
	'P':       (*normal).putHere,
	'R':       (*normal).replaceMode,
	'T':       (*normal).searchCharacterBackwardAfter,
	'V':       (*normal).visualLine,
	'W':       (*normal).wordNonBlank,
	'X':       (*normal).deleteBefore,
	'Y':       (*normal).yankToEnd,
	'a':       (*normal).appendAfter,
	'b':       (*normal).wordBack,
	'c':       (*normal).changeOp,
	'd':       (*normal).deleteOp,
	'e':       (*normal).wordEnd,
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
	's':       (*normal).siegeOp,
	't':       (*normal).searchCharacterBefore,
	'u':       (*normal).undoCmd,
	'v':       (*normal).visual,
	'w':       (*normal).word,
	'x':       (*normal).deleteUnder,
	'y':       (*normal).yankOp,
}

func (e *nvCommon) endline() (_ modeChanger) {
	e.move(len(e.buf))
	return
}

func (e *nvCommon) beginline() (_ modeChanger) {
	e.move(0)
	return
}

func (e *nvCommon) beginlineNonBlank() (_ modeChanger) {
	i := e.indexFunc(isWhitespace, 0, false)
	if i < 0 {
		return e.endline()
	}
	e.move(i)
	return
}

func (e *nvCommon) wordBack() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.wordBackward()
	}
	return
}

func (e *nvCommon) wordBackNonBlank() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.wordBackwardNonBlank()
	}
	return
}

func (e *normal) changeOp() modeChanger {
	return opPend(OpChange, e.count)
}

func (e *normal) deleteOp() modeChanger {
	return opPend(OpDelete, e.count)
}

func (e *normal) yankOp() modeChanger {
	return opPend(OpYank, e.count)
}

func (e *nvCommon) left() (_ modeChanger) {
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

func (e *normal) insertFromBeginning() modeChanger {
	e.move(0)
	return ins(e.pos == len(e.buf))
}

func (e *normal) appendAtEnd() modeChanger {
	e.move(len(e.buf))
	return ins(e.pos == len(e.buf))
}

func (e *normal) insertFirstNonBlank() modeChanger {
	_ = e.beginlineNonBlank()
	return ins(e.pos == len(e.buf))
}

func (e *normal) appendAfter() modeChanger {
	e.move(e.pos + 1)
	return ins(e.pos == len(e.buf))
}

func (e *normal) edit() modeChanger {
	return ins(e.pos == len(e.buf))
}

func (e *normal) down() (_ modeChanger) {
	if e.age >= len(e.history)-e.count {
		return
	}
	e.history[e.age] = e.buf
	e.age += e.count
	e.buf = e.history[e.age]
	e.pos = len(e.buf)
	return
}

func (e *normal) up() (_ modeChanger) {
	if e.age <= e.count-1 {
		return
	}
	if e.age == len(e.history) {
		e.history = append(e.history, e.buf)
	} else {
		e.history[e.age] = e.buf
	}
	e.age -= e.count
	e.buf = e.history[e.age]
	e.pos = len(e.buf)
	return
}

func (e *nvCommon) right() (_ modeChanger) {
	e.move(e.pos + e.count)
	return
}

func (e *normal) putHere() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.put(e.regName, e.pos)
		e.move(e.pos - 1)
	}
	e.undoTree.add(e.buf)
	return
}

func (e *normal) put1() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.put(e.regName, e.pos+1)
		e.move(e.pos + len(e.Read(e.regName)))
	}
	e.undoTree.add(e.buf)
	return
}

func (e *normal) replace() (_ modeChanger) {
	r, _, _ := e.streamSet.in.ReadRune()
	s := make([]rune, e.count)
	for i := 0; i < e.count; i++ {
		s[i] = r
	}
	e.editor.replace(s, e.pos)
	e.move(e.pos + e.count - 1)
	e.undoTree.add(e.buf)
	return
}

func (e *nvCommon) word() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.wordForward()
	}
	return
}

func (e *nvCommon) wordNonBlank() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.wordForwardNonBlank()
	}
	return
}

func (e *normal) wordEnd() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.nvCommon.wordEnd()
	}
	return
}

func (e *normal) wordEndNonBlank() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.nvCommon.wordEndNonBlank()
	}
	return
}

func (e *normal) deleteUnder() (_ modeChanger) {
	from, to := e.pos, e.pos+e.count
	e.yank(e.regName, from, to)
	e.delete(from, to)
	e.undoTree.add(e.buf)
	return
}

func (e *normal) deleteBefore() (_ modeChanger) {
	from, to := e.pos-e.count, e.pos
	e.yank(e.regName, from, to)
	e.delete(from, to)
	e.undoTree.add(e.buf)
	return
}

func (e *normal) deleteToEnd() (_ modeChanger) {
	from, to := e.pos, len(e.buf)
	e.yank(e.regName, from, to)
	e.delete(from, to)
	e.undoTree.add(e.buf)
	return
}

func (e *normal) changeToEnd() modeChanger {
	_ = e.deleteToEnd()
	return ins(e.pos == len(e.buf))
}

func (e *normal) yankToEnd() (_ modeChanger) {
	e.yank(e.regName, e.pos, len(e.buf))
	return
}

func (e *nvCommon) searchCharacter() (_ modeChanger) {
	r, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return
	}
	pos := e.pos
	for i := 0; i < e.count; i++ {
		i, err := e.charSearch(r)
		if err != nil {
			e.move(pos)
			return
		}
		e.move(i)
	}
	return
}

func (e *nvCommon) searchCharacterBackward() (_ modeChanger) {
	r, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return
	}
	pos := e.pos
	for i := 0; i < e.count; i++ {
		i, err := e.charSearchBackward(r)
		if err != nil {
			e.move(pos)
			return
		}
		e.move(i)
	}
	return
}

func (e *nvCommon) searchCharacterBefore() (_ modeChanger) {
	r, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return
	}
	pos := e.pos
	for i := 0; i < e.count; i++ {
		i, err := e.charSearchBefore(r)
		if err != nil {
			e.move(pos)
			return
		}
		e.move(i + 1)
	}
	e.move(e.pos - 1)
	return
}

func (e *nvCommon) searchCharacterBackwardAfter() (_ modeChanger) {
	r, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return
	}
	pos := e.pos
	for i := 0; i < e.count; i++ {
		i, err := e.charSearchBackwardAfter(r)
		if err != nil {
			e.move(pos)
			return
		}
		e.move(i - 1)
	}
	e.move(e.pos + 1)
	return
}

func (e *normal) gCmd() (_ modeChanger) {
	r, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return
	}
	switch r {
	case 'u':
		return opPend(OpLower, e.count)
	case 'U':
		return opPend(OpUpper, e.count)
	case '~':
		return opPend(OpTilde, e.count)
	case '/':
		return e.searchHistory()
	case 'I':
		return e.insertFromBeginning()
	case 'e':
		return e.wordEndBackward()
	case 'E':
		return e.wordEndBackwardNonBlank()
	}
	return
}

func (e *normal) undoCmd() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.undo()
	}
	return
}

func (e *normal) redoCmd() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.redo()
	}
	return
}

func (e *normal) replaceMode() modeChanger {
	return func(b *balancer) (moder, error) {
		return newReplace(b.streamSet, b.editor), nil
	}
}

func (e *normal) handleRegister() (_ modeChanger) {
	r, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return
	}
	if !register.IsValid(r) {
		return
	}
	e.regName = r
	return
}

func (e *normal) commandline() modeChanger {
	return func(b *balancer) (moder, error) {
		return newCommandline(b.streamSet, b.editor), nil
	}
}

func (e *normal) visual() modeChanger {
	return func(b *balancer) (moder, error) {
		return newVisual(b.streamSet, b.editor), nil
	}
}

func (e *normal) visualLine() modeChanger {
	return func(b *balancer) (moder, error) {
		return newVisualLine(b.streamSet, b.editor), nil
	}
}

func (e *nvCommon) column() (_ modeChanger) {
	e.move(constrain(e.count-1, 0, len(e.buf)))
	return
}

func (e *normal) search() modeChanger {
	return func(b *balancer) (moder, error) {
		return newSearch(b.streamSet, b.editor, searchForward), nil
	}
}

func (e *normal) searchBackward() modeChanger {
	return func(b *balancer) (moder, error) {
		return newSearch(b.streamSet, b.editor, searchBackward), nil
	}
}

func (e *normal) next() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.move(e.nvCommon.next())
	}
	return
}

func (e *normal) previous() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.move(e.nvCommon.previous())
	}
	return
}

func (e *normal) searchHistory() (_ modeChanger) {
	return func(b *balancer) (moder, error) {
		return newSearch(b.streamSet, b.editor, searchHistoryForward), nil
	}
}

func (e *normal) indexNumber() int {
	i0 := e.pos
	for ; 0 <= i0; i0-- {
		r := e.buf[i0]
		if !('0' <= r && r <= '9') {
			break
		}
	}
	var negative bool
	for i, r := range e.buf[i0:] {
		if '0' <= r && r <= '9' {
			if negative {
				i--
			}
			return i + i0
		}
		if !negative && r == '-' {
			negative = true
		}
	}
	return -1
}

func (e *normal) parseNumber(i int) (a int, l int) {
	i0 := i
	if e.buf[i] == '-' {
		i0++
		defer func() {
			a *= -1
		}()
	}
	for n := i0; n < len(e.buf); n++ {
		r := int(e.buf[n])
		if '0' <= r && r <= '9' {
			a = 10*a + r - '0'
			continue
		}
		return a, n - i
	}
	return a, len(e.buf) - i
}

func (e *normal) updateNumber(f func(int) int) {
	i := e.indexNumber()
	if i < 0 {
		return
	}
	n, l := e.parseNumber(i)
	e.delete(i, i+l)
	s := fmt.Sprint(f(n))
	e.insert([]rune(s), i)
	e.move(i + len(s) - 1)
}

func (e *normal) increment() (_ modeChanger) {
	e.updateNumber(func(n int) int {
		return n + e.count
	})
	return
}

func (e *normal) decrement() (_ modeChanger) {
	e.updateNumber(func(n int) int {
		return n - e.count
	})
	return
}

func (e *normal) switchCase() (_ modeChanger) {
	e.editor.swapCase(e.pos, e.pos+e.count)
	e.move(e.pos + e.count)
	return
}

func (e *normal) siegeOp() (_ modeChanger) {
	return opPend(OpSiege, e.count)
}

func (e *nvCommon) prevUnmatched() (_ modeChanger) {
	r, _, _ := e.in.ReadRune()
	var rp, lp rune
	switch r {
	case '(':
		rp = '('
		lp = ')'
	case '{':
		rp = '{'
		lp = '}'
	default:
		return
	}
	i := e.searchLeft(rp, lp)
	if i < 0 {
		return
	}
	e.move(i)
	return
}

func (e *nvCommon) nextUnmatched() (_ modeChanger) {
	r, _, _ := e.in.ReadRune()
	var rp, lp rune
	switch r {
	case ')':
		rp = '('
		lp = ')'
	case '}':
		rp = '{'
		lp = '}'
	default:
		return
	}
	i := e.searchRight(rp, lp)
	if i < 0 {
		return
	}
	e.move(i)
	return
}

func isParen(r rune) bool {
	switch r {
	case '(', ')', '[', ']', '{', '}':
		return true
	}
	return false
}

func getRightParen(r rune) rune {
	return map[rune]rune{
		'(': ')',
		'[': ']',
		'{': '}',
	}[r]
}

func getLeftParen(r rune) rune {
	return map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}[r]
}

func (e *nvCommon) moveToMatch() (_ modeChanger) {
	i := e.indexFunc(isParen, e.pos, true)
	if i < 0 {
		return
	}
	initPos := e.pos
	e.move(i)
	p := e.buf[i]
	switch p {
	case '(', '[', '{':
		i := e.searchRight(p, getRightParen(p))
		if i < 0 {
			e.move(initPos)
			return
		}
		e.move(i)
	case ')', ']', '}':
		i := e.searchLeft(getLeftParen(p), p)
		if i < 0 {
			e.move(initPos)
			return
		}
		e.move(i)
	}
	return
}

func (e *normal) wordEndBackwardNonBlank() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.editor.wordEndBackwardNonBlank()
	}
	return
}

func (e *normal) wordEndBackward() (_ modeChanger) {
	for i := 0; i < e.count; i++ {
		e.editor.wordEndBackward()
	}
	return
}

type operatorPending struct {
	nvCommon

	opType     int
	opCount    int
	start      int
	inclusive  bool
	motionType int
}

func newOperatorPending(s streamSet, e *editor, op int, count int) *operatorPending {
	return &operatorPending{
		nvCommon: nvCommon{
			streamSet: s,
			editor:    e,
		},
		opType:  op,
		opCount: count,
		start:   e.pos,
	}
}

func opPend(op, count int) modeChanger {
	return func(b *balancer) (moder, error) {
		return newOperatorPending(b.streamSet, b.editor, op, count), nil
	}
}

func (o *operatorPending) Mode() mode {
	return modeOperatorPending
}

func (o *operatorPending) Runes() []rune {
	return o.buf
}

func (o *operatorPending) Position() int {
	return o.pos
}

func (o *operatorPending) Message() []rune {
	return nil
}

func (o *operatorPending) Highlight() *screen.Hi {
	return nil
}

func (o *operatorPending) Run() (end continuity, next modeChanger, err error) {
	next = norm()
	r, _, err := o.in.ReadRune()
	if err != nil {
		return end, next, err
	}
	for ('1' <= r && r <= '9') || (o.count != 0 && r == '0') {
		o.count = o.count*10 + int(r-'0')
		r1, _, err := o.streamSet.in.ReadRune()
		if err != nil {
			return end, next, err
		}
		r = r1
	}
	if o.count == 0 {
		o.count = 1
	}
	if o.opCount > 0 {
		o.count *= o.opCount
	}
	cmd, ok := operatorPendingCommands[r]
	if !ok {
		return
	}
	if m := cmd(o); m != nil {
		next = m
	}
	o.count = 0

	if m := o.operate(); m != nil {
		next = m
	}

	if o.pos == len(o.buf) {
		o.move(o.pos - 1)
	}

	return
}

type operatorPendingCommand = func(*operatorPending) modeChanger

var operatorPendingCommands = map[rune]operatorPendingCommand{
	'$': (*operatorPending).endline,
	'%': (*operatorPending).moveToMatch,
	'[': (*operatorPending).prevUnmatched,
	']': (*operatorPending).nextUnmatched,
	'|': (*operatorPending).column,
	'^': (*operatorPending).beginlineNonBlank,
	'0': (*operatorPending).beginline,
	'B': (*operatorPending).wordBackNonBlank,
	'E': (*operatorPending).wordEndNonBlank,
	'F': (*operatorPending).searchCharacterBackward,
	'T': (*operatorPending).searchCharacterBackwardAfter,
	'W': (*operatorPending).wordNonBlank,
	'a': (*operatorPending).anObject,
	'b': (*operatorPending).wordBack,
	'c': (*operatorPending).changeLine,
	'd': (*operatorPending).deleteLine,
	'e': (*operatorPending).wordEnd,
	'f': (*operatorPending).searchCharacter,
	'h': (*operatorPending).left,
	'i': (*operatorPending).innerObject,
	'l': (*operatorPending).right,
	't': (*operatorPending).searchCharacterBefore,
	'w': (*operatorPending).word,
	'y': (*operatorPending).yankLine,
}

func (o *operatorPending) operate() modeChanger {
	from := min(o.start, o.pos)
	to := max(o.start, o.pos)
	if o.inclusive {
		to++
	}
	if o.motionType == mline {
		from = 0
		to = len(o.buf)
	}
	switch o.opType {
	case OpDelete:
		// FIXME: specify register
		o.yank(register.Unnamed, from, to)
		o.delete(from, to)
		o.undoTree.add(o.buf)
	case OpYank:
		o.yank(register.Unnamed, from, to)
		if o.motionType == mline {
			// yanking line does not move cursor.
			return nil
		}
	case OpChange:
		o.yank(register.Unnamed, from, to)
		o.delete(from, to)
		o.undoTree.add(o.buf)
		return ins(o.pos == len(o.buf))
	case OpLower:
		o.toLower(from, to)
		o.undoTree.add(o.buf)
	case OpUpper:
		o.toUpper(from, to)
		o.undoTree.add(o.buf)
	case OpTilde:
		o.swapCase(from, to)
		o.undoTree.add(o.buf)
	case OpSiege:
		r, _, _ := o.in.ReadRune()
		o.siege(from, to, r)
	}
	o.move(min(from, to))
	return nil
}

func (o *operatorPending) wordEnd() (_ modeChanger) {
	for i := 0; i < o.count; i++ {
		o.nvCommon.wordEnd()
	}
	o.inclusive = true
	return
}

func (o *operatorPending) wordEndNonBlank() (_ modeChanger) {
	for i := 0; i < o.count; i++ {
		o.nvCommon.wordEndNonBlank()
	}
	o.inclusive = true
	return
}

func (o *operatorPending) object(include bool) {
	var from, to int
	r, _, _ := o.in.ReadRune()
	switch r {
	case 'w':
		from, to = o.currentWord(include)
	case 'W':
		from, to = o.currentWordNonBlank(include)
	case '"', '\'', '`':
		from, to = o.currentQuote(include, r)
	case '(', ')', 'b':
		from, to = o.currentParen(include, '(', ')')
	case '{', '}', 'B':
		from, to = o.currentParen(include, '{', '}')
	case '[', ']':
		from, to = o.currentParen(include, '[', ']')
	case '<', '>':
		from, to = o.currentParen(include, '<', '>')
	default:
		return
	}
	if from < 0 || to < 0 {
		return
	}
	o.start = from
	o.pos = to
}

func (o *operatorPending) innerObject() (_ modeChanger) {
	o.object(false)
	return
}

func (o *operatorPending) anObject() (_ modeChanger) {
	o.object(true)
	return
}

func (o *operatorPending) linewise(op int) {
	if o.opType != op {
		return
	}
	o.motionType = mline
}

func (o *operatorPending) deleteLine() (_ modeChanger) {
	o.linewise(OpDelete)
	return
}

func (o *operatorPending) changeLine() (_ modeChanger) {
	o.linewise(OpChange)
	return
}

func (o *operatorPending) yankLine() (_ modeChanger) {
	o.linewise(OpYank)
	return
}
