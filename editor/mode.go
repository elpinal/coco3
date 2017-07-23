package editor

type mode int

const (
	modeNormal mode = iota + 1
	modeVisual
	modeSelect
	modeInsert
	modeCommandline
	modeEx

	modeOperatorPending
	modeReplace
	modeVirtualReplace
	modeInsertNormal
	modeInsertVisual
	modeInsertSelect
)

type moder interface {
	Mode() mode
	Run() (end continuity, next mode, err error)
	Runes() []rune
	Position() int
	Message() []rune
}
