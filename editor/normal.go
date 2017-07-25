package editor

import "github.com/elpinal/coco3/editor/register"

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
			editor: e,
		},
	}
}

func (e *normal) Mode() mode {
	return modeNormal
}

func (e *normal) Run() (end continuity, next mode, err error) {
	next = modeNormal
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
		if m := cmd(e, r); m != 0 {
			next = m
		}
		if n := e.doPendingOperator(); n != 0 {
			next = n
		}
		e.count = 0
		if next != modeInsert && e.pos == len(e.buf) {
			e.move(e.pos - 1)
		}
	}
	return end, next, err
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

type normalCommand = func(*normal, rune) mode

var normalCommands = map[rune]normalCommand{
	CharCtrlR: (*normal).redoCmd,
	'"':       (*normal).handleRegister,
	':':       (*normal).commandline,
	'$':       (*normal).endline,
	'0':       (*normal).beginline,
	'A':       (*normal).edit,
	'B':       (*normal).wordBack,
	'C':       (*normal).abbrev,
	'D':       (*normal).abbrev,
	'E':       (*normal).word,
	'F':       (*normal).searchBackward,
	'I':       (*normal).edit,
	'R':       (*normal).replaceMode,
	'W':       (*normal).word,
	'X':       (*normal).abbrev,
	'Y':       (*normal).abbrev,
	'a':       (*normal).edit,
	'b':       (*normal).wordBack,
	'c':       (*normal).operator1,
	'd':       (*normal).operator1,
	'e':       (*normal).word,
	'f':       (*normal).search,
	'g':       (*normal).gCmd,
	'h':       (*normal).left,
	'i':       (*normal).edit,
	'j':       (*normal).down,
	'k':       (*normal).up,
	'l':       (*normal).right,
	'p':       (*normal).put1,
	'r':       (*normal).replace,
	'u':       (*normal).undoCmd,
	'v':       (*normal).visual,
	'w':       (*normal).word,
	'x':       (*normal).abbrev,
	'y':       (*normal).operator1,
}

func (e *nvCommon) endline(r rune) (next mode) {
	e.move(len(e.buf))
	return
}

func (e *nvCommon) beginline(r rune) (next mode) {
	e.move(0)
	return
}

func (e *nvCommon) wordBack(r rune) (next mode) {
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

func (e *normal) operator1(r rune) mode {
	return e.operator(string(r))
}

func (e *normal) operator(s string) mode {
	op := opChars[s]
	if op == e.opType { // double operator
		e.motionType = mline
	} else {
		e.opStart = e.pos
		e.opType = op
		e.opCount = e.count
	}
	return modeNormal
}

func (e *nvCommon) left(r rune) (next mode) {
	e.move(e.pos - e.count)
	return
}

func (e *normal) edit(r rune) mode {
	if (r == 'a' || r == 'i') && e.opType != OpNop {
		e.object(r)
		return modeNormal
	}
	switch r {
	case 'A':
		e.move(len(e.buf))
	case 'I':
		e.move(0)
	case 'a':
		e.move(e.pos + 1)
	}
	return modeInsert
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

func (e *normal) down(r rune) mode {
	if e.age >= len(e.history)-1 {
		return modeNormal
	}
	e.history[e.age] = e.buf
	e.age++
	e.buf = e.history[e.age]
	e.pos = len(e.buf)
	return modeNormal
}

func (e *normal) up(r rune) mode {
	if e.age <= 0 {
		return modeNormal
	}
	if e.age == len(e.history) {
		e.history = append(e.history, e.buf)
	} else {
		e.history[e.age] = e.buf
	}
	e.age--
	e.buf = e.history[e.age]
	e.pos = len(e.buf)
	return modeNormal
}

func (e *nvCommon) right(r rune) (next mode) {
	e.move(e.pos + e.count)
	return
}

func (e *normal) put1(r rune) mode {
	for i := 0; i < e.count; i++ {
		e.put(e.regName, e.pos+1)
	}
	e.undoTree.add(e.buf)
	return modeNormal
}

func (e *normal) replace(r rune) mode {
	r1, _, _ := e.streamSet.in.ReadRune()
	s := make([]rune, e.count)
	for i := 0; i < e.count; i++ {
		s[i] = r1
	}
	e.editor.replace(s, e.pos)
	e.move(e.pos + e.count - 1)
	e.undoTree.add(e.buf)
	return modeNormal
}

func (e *normal) word(r rune) mode {
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
	return modeNormal
}

func (e *normal) doPendingOperator() mode {
	if !e.finishOp {
		return 0
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
		return modeInsert
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
	return modeNormal
}

func (e *normal) clearOp() {
	e.opArg = opArg{}
}

func (e *normal) abbrev(r rune) mode {
	amap := map[rune][]rune{
		'x': []rune("dl"),
		'X': []rune("dh"),
		'D': []rune("d$"),
		'C': []rune("c$"),
		'Y': []rune("y$"),
	}
	e.streamSet.in.Add(amap[r])
	e.opCount = e.count
	return modeNormal
}

func (e *nvCommon) search(r rune) (next mode) {
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

func (e *nvCommon) searchBackward(r rune) (next mode) {
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

func (e *normal) gCmd(r rune) mode {
	r1, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return modeNormal
	}
	switch r1 {
	case 'u', 'U', '~':
		return e.operator(string([]rune{r, r1}))
	}
	return modeNormal
}

func (e *normal) undoCmd(r rune) mode {
	e.undo()
	return modeNormal
}

func (e *normal) redoCmd(r rune) mode {
	e.redo()
	return modeNormal
}

func (e *normal) replaceMode(r rune) mode {
	return modeReplace
}

func (e *normal) handleRegister(r rune) mode {
	r1, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return modeNormal
	}
	if !register.IsValid(r1) {
		return modeNormal
	}
	e.regName = r1
	return modeNormal
}

func (e *normal) commandline(r rune) mode {
	return modeCommandline
}

func (e *normal) visual(r rune) mode {
	return modeVisual
}
