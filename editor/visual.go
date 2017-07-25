package editor

type visual struct {
	nvCommon
}

func newVisual(s streamSet, e *editor) *visual {
	return &visual{
		nvCommon: nvCommon{
			streamSet: s,
			editor:    e,
		},
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

func (v *visual) Run() (end continuity, next mode, err error) {
	next = modeVisual
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
	if m := cmd(v, r); m != 0 {
		next = m
	}
	if next != modeInsert && v.pos == len(v.buf) {
		v.move(v.pos - 1)
	}
	v.count = 0
	return
}

type visualCommand = func(*visual, rune) mode

var visualCommands = map[rune]visualCommand{
	CharEscape: (*visual).escape,
	'$':        (*visual).endline,
	'0':        (*visual).beginline,
	'B':        (*visual).wordBack,
	'F':        (*visual).searchBackward,
	'b':        (*visual).wordBack,
	'f':        (*visual).search,
	'h':        (*visual).left,
	'l':        (*visual).right,
}

func (v *visual) escape(_ rune) mode {
	return modeNormal
}
