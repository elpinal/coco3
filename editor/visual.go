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
	cmd, ok := visualCommands[r]
	if !ok {
		return
	}
	next = cmd(v, r)
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
