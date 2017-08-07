package editor

import "github.com/elpinal/coco3/screen"

const (
	OpNop = iota
	OpDelete
	OpYank
	OpChange
	OpLower
	OpUpper
	OpSwitchCase
	OpSiege
)

type operatorPending struct {
	nvCommon

	opType     int
	opCount    int
	start      int
	inclusive  bool
	motionType int
	regName    rune
}

func newOperatorPending(s streamSet, e *editor, op, count int, regName rune) *operatorPending {
	return &operatorPending{
		nvCommon: nvCommon{
			streamSet: s,
			editor:    e,
		},
		opType:  op,
		opCount: count,
		start:   e.pos,
		regName: regName,
	}
}

func opPend(op, count int, regName rune) modeChanger {
	return func(b *balancer) (moder, error) {
		return newOperatorPending(b.streamSet, b.editor, op, count, regName), nil
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
	'~': (*operatorPending).switchLine,
	'0': (*operatorPending).beginline,
	'B': (*operatorPending).wordBackNonBlank,
	'E': (*operatorPending).wordEndNonBlank,
	'F': (*operatorPending).searchCharacterBackward,
	'N': (*operatorPending).previous,
	'T': (*operatorPending).searchCharacterBackwardAfter,
	'U': (*operatorPending).upperLine,
	'W': (*operatorPending).wordNonBlank,
	'a': (*operatorPending).anObject,
	'b': (*operatorPending).wordBack,
	'c': (*operatorPending).changeLine,
	'd': (*operatorPending).deleteLine,
	'e': (*operatorPending).wordEnd,
	'f': (*operatorPending).searchCharacter,
	'g': (*operatorPending).gCmd,
	'h': (*operatorPending).left,
	'i': (*operatorPending).innerObject,
	'l': (*operatorPending).right,
	'n': (*operatorPending).next,
	't': (*operatorPending).searchCharacterBefore,
	'u': (*operatorPending).lowerLine,
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
		o.yank(o.regName, from, to)
		o.delete(from, to)
		// TODO: It is hard to remember to write "undoTree.add" every time
		// changing text.
		o.undoTree.add(o.buf)
	case OpYank:
		o.yank(o.regName, from, to)
		if o.motionType == mline {
			// yanking line does not move cursor.
			return nil
		}
	case OpChange:
		o.yank(o.regName, from, to)
		o.delete(from, to)
		o.undoTree.add(o.buf)
		return ins(o.pos == len(o.buf))
	case OpLower:
		o.toLower(from, to)
		o.undoTree.add(o.buf)
	case OpUpper:
		o.toUpper(from, to)
		o.undoTree.add(o.buf)
	case OpSwitchCase:
		o.switchCase(from, to)
		o.undoTree.add(o.buf)
	case OpSiege:
		r, _, _ := o.in.ReadRune()
		o.siege(from, to, r)
		o.undoTree.add(o.buf)
	}
	o.move(min(from, to))
	return nil
}

func (o *operatorPending) wordEnd() (_ modeChanger) {
	_ = o.nvCommon.wordEnd()
	o.inclusive = true
	return
}

func (o *operatorPending) wordEndNonBlank() (_ modeChanger) {
	_ = o.nvCommon.wordEndNonBlank()
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

func (o *operatorPending) switchLine() (_ modeChanger) {
	o.linewise(OpSwitchCase)
	return
}

func (o *operatorPending) upperLine() (_ modeChanger) {
	o.linewise(OpUpper)
	return
}

func (o *operatorPending) lowerLine() (_ modeChanger) {
	o.linewise(OpLower)
	return
}

func (o *operatorPending) gCmd() (_ modeChanger) {
	r, _, err := o.in.ReadRune()
	if err != nil {
		return
	}
	switch r {
	case 'u':
		o.linewise(OpLower)
	case 'U':
		o.linewise(OpUpper)
	case '~':
		o.linewise(OpSwitchCase)
	case 'e':
		o.inclusive = true
		return o.nvCommon.wordEndBackward()
	case 'E':
		o.inclusive = true
		return o.nvCommon.wordEndBackwardNonBlank()
	}
	return
}
