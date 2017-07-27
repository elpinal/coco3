package editor

import (
	"fmt"

	"github.com/elpinal/coco3/screen"
)

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

	// Not in Vim
	modeSearch
)

type moder interface {
	Mode() mode
	Run() (end continuity, next mode, err error)
	Runes() []rune
	Position() int
	Message() []rune
	Highlight() *screen.Hi
}

func (m mode) String() string {
	switch m {
	case modeNormal:
		return "normal"
	case modeVisual:
		return "visual"
	case modeSelect:
		return "select"
	case modeInsert:
		return "insert"
	case modeCommandline:
		return "command-line"
	case modeEx:
		return "ex"
	case modeOperatorPending:
		return "operator-pending"
	case modeReplace:
		return "replace"
	case modeVirtualReplace:
		return "virtual replace"
	case modeInsertNormal:
		return "insert normal"
	case modeInsertVisual:
		return "insert visual"
	case modeInsertSelect:
		return "insert select"
	case modeSearch:
		return "search"
	}
	return fmt.Sprintf("number (%d)", m)
}
