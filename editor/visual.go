package editor

import (
	"github.com/elpinal/coco3/editor/register"
	"github.com/elpinal/coco3/screen"
)

type visual struct {
	nvCommon

	start int
}

func newVisual(s streamSet, e *editor) *visual {
	return &visual{
		nvCommon: nvCommon{
			streamSet: s,
			editor:    e,
		},

		start: e.pos,
	}
}

func (v *visual) Mode() mode {
	return modeVisual
}

func (v *visual) Runes() []rune {
	return v.buf
}

func (v *visual) Position() int {
	return v.pos
}

func (v *visual) Message() []rune {
	return []rune("-- VISUAL --")
}

func (v *visual) Highlight() *screen.Hi {
	return &screen.Hi{
		Left:  min(v.start, v.pos),
		Right: constrain(max(v.start, v.pos)+1, 0, len(v.buf)),
	}
}

func (v *visual) Run() (end continuity, next modeChanger, err error) {
	r, _, err := v.in.ReadRune()
	if err != nil {
		return end, next, err
	}
	for ('1' <= r && r <= '9') || (v.count != 0 && r == '0') {
		v.count = v.count*10 + int(r-'0')
		r1, _, err := v.streamSet.in.ReadRune()
		if err != nil {
			return end, next, err
		}
		r = r1
	}
	if v.count == 0 {
		v.count = 1
	}
	cmd, ok := visualCommands[r]
	if !ok {
		return
	}
	if m := cmd(v, r); m != nil {
		next = m
	}
	if v.pos == len(v.buf) {
		v.move(v.pos - 1)
	}
	v.count = 0
	return
}

type visualCommand = func(*visual, rune) modeChanger

var visualCommands = map[rune]visualCommand{
	CharEscape: (*visual).escape,
	CharCtrlC:  (*visual).escape,
	'$':        (*visual).endline,
	'0':        (*visual).beginline,
	'B':        (*visual).wordBack,
	'E':        (*visual).wordEnd,
	'F':        (*visual).searchCharacterBackward,
	'U':        (*visual).toUpper,
	'W':        (*visual).wordNonBlank,
	'b':        (*visual).wordBack,
	'c':        (*visual).change,
	'd':        (*visual).delete,
	'e':        (*visual).wordEnd,
	'f':        (*visual).searchCharacter,
	'h':        (*visual).left,
	'l':        (*visual).right,
	'r':        (*visual).replace,
	'o':        (*visual).swap,
	'u':        (*visual).toLower,
	'v':        (*visual).escape,
	'w':        (*visual).word,
	'y':        (*visual).yank,
}

func (v *visual) escape(_ rune) modeChanger {
	return norm()
}

func (v *visual) delete(_ rune) modeChanger {
	hi := v.Highlight()
	v.nvCommon.delete(hi.Left, hi.Right)
	return norm()
}

func (v *visual) change(_ rune) modeChanger {
	_ = v.delete(0)
	return ins(v.pos == len(v.buf))
}

func (v *visual) swap(_ rune) (_ modeChanger) {
	v.start, v.pos = v.pos, v.start
	return
}

func (v *visual) yank(_ rune) modeChanger {
	hi := v.Highlight()
	v.nvCommon.yank(register.Unnamed, hi.Left, hi.Right)
	v.move(hi.Left)
	return norm()
}

func (v *visual) replace(_ rune) modeChanger {
	r, _, _ := v.in.ReadRune()
	hi := v.Highlight()
	rs := make([]rune, hi.Right-hi.Left)
	for i := range rs {
		rs[i] = r
	}
	v.nvCommon.replace(rs, hi.Left)
	return norm()
}

func (v *visual) toUpper(_ rune) modeChanger {
	hi := v.Highlight()
	v.nvCommon.toUpper(hi.Left, hi.Right)
	return norm()
}

func (v *visual) toLower(_ rune) modeChanger {
	hi := v.Highlight()
	v.nvCommon.toLower(hi.Left, hi.Right)
	return norm()
}

func (v *visual) wordEnd(r rune) (_ modeChanger) {
	var f func()
	switch r {
	case 'e':
		f = v.nvCommon.wordEnd
	case 'E':
		f = v.nvCommon.wordEndNonBlank
	}
	for i := 0; i < v.count; i++ {
		f()
	}
	return
}
