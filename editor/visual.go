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
	if m := cmd(v); m != nil {
		next = m
	}
	if v.pos == len(v.buf) {
		v.move(v.pos - 1)
	}
	v.count = 0
	return
}

type visualCommand = func(*visual) modeChanger

var visualCommands = map[rune]visualCommand{
	CharEscape: (*visual).escape,
	CharCtrlC:  (*visual).escape,
	'~':        (*visual).switchCase,
	'$':        (*visual).endline,
	'0':        (*visual).beginline,
	'A':        (*visual).appendAfter,
	'B':        (*visual).wordBack,
	'C':        (*visual).changeLine,
	'D':        (*visual).deleteLine,
	'E':        (*visual).wordEndNonBlank,
	'F':        (*visual).searchCharacterBackward,
	'I':        (*visual).insertBefore,
	'U':        (*visual).toUpper,
	'W':        (*visual).wordNonBlank,
	'b':        (*visual).wordBack,
	'c':        (*visual).change,
	'd':        (*visual).delete,
	'e':        (*visual).wordEnd,
	'f':        (*visual).searchCharacter,
	'i':        (*visual).object,
	'h':        (*visual).left,
	'l':        (*visual).right,
	'r':        (*visual).replace,
	's':        (*visual).siege,
	'o':        (*visual).swap,
	'u':        (*visual).toLower,
	'v':        (*visual).escape,
	'w':        (*visual).word,
	'y':        (*visual).yank,
}

func (v *visual) escape() modeChanger {
	return norm()
}

func (v *visual) delete() modeChanger {
	hi := v.Highlight()
	v.nvCommon.delete(hi.Left, hi.Right)
	return norm()
}

func (v *visual) change() modeChanger {
	_ = v.delete()
	return ins(v.pos == len(v.buf))
}

func (v *visual) swap() (_ modeChanger) {
	v.start, v.pos = v.pos, v.start
	return
}

func (v *visual) yank() modeChanger {
	hi := v.Highlight()
	v.nvCommon.yank(register.Unnamed, hi.Left, hi.Right)
	v.move(hi.Left)
	return norm()
}

func (v *visual) replace() modeChanger {
	r, _, _ := v.in.ReadRune()
	hi := v.Highlight()
	rs := make([]rune, hi.Right-hi.Left)
	for i := range rs {
		rs[i] = r
	}
	v.nvCommon.replace(rs, hi.Left)
	return norm()
}

func (v *visual) toUpper() modeChanger {
	hi := v.Highlight()
	v.nvCommon.toUpper(hi.Left, hi.Right)
	return norm()
}

func (v *visual) toLower() modeChanger {
	hi := v.Highlight()
	v.nvCommon.toLower(hi.Left, hi.Right)
	return norm()
}

func (v *visual) wordEnd() (_ modeChanger) {
	for i := 0; i < v.count; i++ {
		v.nvCommon.wordEnd()
	}
	return
}

func (v *visual) wordEndNonBlank() (_ modeChanger) {
	for i := 0; i < v.count; i++ {
		v.nvCommon.wordEndNonBlank()
	}
	return
}

func (v *visual) insertBefore() modeChanger {
	v.move(v.Highlight().Left)
	return ins(v.pos == len(v.buf))
}

func (v *visual) appendAfter() modeChanger {
	v.move(v.Highlight().Right)
	return ins(v.pos == len(v.buf))
}

func (v *visual) switchCase() modeChanger {
	hi := v.Highlight()
	v.editor.swapCase(hi.Left, hi.Right)
	v.move(hi.Left)
	return norm()
}

func (v *visual) deleteLine() modeChanger {
	v.editor.delete(0, len(v.buf))
	return norm()
}

func (v *visual) siege() modeChanger {
	hi := v.Highlight()
	r, _, _ := v.in.ReadRune()
	v.editor.siege(hi.Left, hi.Right, r)
	v.move(hi.Left)
	return norm()
}

func (v *visual) changeLine() modeChanger {
	v.editor.delete(0, len(v.buf))
	return ins(v.pos == len(v.buf))
}

func (v *visual) object() (_ modeChanger) {
	if v.start <= v.pos {
		v.move(v.pos + 1)
	} else {
		v.move(v.pos - 1)
	}
	from, to := v.currentWord(false)
	if from < 0 {
		return
	}
	if v.start <= v.pos {
		v.move(to - 1)
		return
	}
	v.move(from)
	return
}
