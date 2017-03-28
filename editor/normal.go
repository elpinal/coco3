package editor

import "github.com/elpinal/coco3/editor/register"

type opArg struct {
	opType     int // current operator type
	opStart    int
	opCount    int
	motionType int
}

type normalSet struct {
	opArg
	finishOp bool
	count    int
}

type normal struct {
	streamSet
	*editor

	normalSet
}

func (e *normal) Mode() mode {
	return modeNormal
}

func (e *normal) Run() (end bool, next mode, err error) {
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
	for _, cmd := range normalCommands {
		if cmd.r == r {
			next = cmd.fn(e, r)
			if n := e.doPendingOperator(); n != 0 {
				next = n
			}
			e.count = 0
			if next != modeInsert && e.pos == len(e.buf) {
				e.move(e.pos - 1)
			}
			return
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

type normalCommand struct {
	r   rune                     // first char
	fn  func(*normal, rune) mode // function for this command
	arg int
}

var normalCommands = []normalCommand{
	{'$', (*normal).endline, 0},
	{'0', (*normal).beginline, 0},
	{'A', (*normal).edit, 0},
	{'B', (*normal).wordBack, 0},
	{'C', (*normal).abbrev, 0},
	{'D', (*normal).abbrev, 0},
	{'I', (*normal).edit, 0},
	{'W', (*normal).word, 0},
	{'X', (*normal).abbrev, 0},
	{'Y', (*normal).abbrev, 0},
	{'a', (*normal).edit, 0},
	{'b', (*normal).wordBack, 0},
	{'c', (*normal).operator, 0},
	{'d', (*normal).operator, 0},
	{'e', (*normal).word, 0},
	{'h', (*normal).left, 0},
	{'i', (*normal).edit, 0},
	{'j', (*normal).down, 0},
	{'k', (*normal).up, 0},
	{'l', (*normal).right, 0},
	{'p', (*normal).put1, 0},
	{'r', (*normal).replace, 0},
	{'w', (*normal).word, 0},
	{'x', (*normal).abbrev, 0},
	{'y', (*normal).operator, 0},
}

func (e *normal) endline(r rune) mode {
	e.move(len(e.buf))
	return modeNormal
}

func (e *normal) beginline(r rune) mode {
	e.move(0)
	return modeNormal
}

func (e *normal) wordBack(r rune) mode {
	for i := 0; i < e.count; i++ {
		switch r {
		case 'b':
			e.wordBackward()
		case 'B':
			e.wordBackwardNonBlank()
		}
	}
	return modeNormal
}

func (e *normal) operator(r rune) mode {
	op := opChars[r]
	if op == e.opType { // double operator
		e.motionType = mline
	} else {
		e.opStart = e.pos
		e.opType = op
		e.opCount = e.count
	}
	return modeNormal
}

func (e *normal) left(r rune) mode {
	e.move(e.pos - e.count)
	return modeNormal
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

func (e *normal) right(r rune) mode {
	e.move(e.pos + e.count)
	return modeNormal
}

func (e *normal) put1(r rune) mode {
	for i := 0; i < e.count; i++ {
		e.put(register.Unnamed, e.pos+1)
	}
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
			e.wordEnd()
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
	if e.motionType == mline {
		from = 0
		to = len(e.buf)
	}
	switch e.opType {
	case OpDelete:
		e.yank(register.Unnamed, from, to)
		e.delete(from, to)
	case OpYank:
		e.yank(register.Unnamed, from, to)
	case OpChange:
		e.yank(register.Unnamed, from, to)
		e.delete(from, to)
		return modeInsert
	}
	e.clearOp()
	return modeNormal
}

func (e *normal) clearOp() {
	e.opType = OpNop
	e.opCount = 0
	e.motionType = mchar
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
