package editor

import (
	"fmt"
	"strings"

	"github.com/elpinal/coco3/screen"
)

type exCommand struct {
	name string
	fn   func(*commandline, []string) continuity
}

// exComands represents a table of Ex commands and corresponding functions.
// The order is important. Precede commands have higher precedence.
var exCommands = []exCommand{
	{"help", (*commandline).help},
	{"delete", (*commandline).delete},
	{"quit", (*commandline).quit},
}

type commandline struct {
	streamSet
	*editor

	basic *basic
}

func newCommandline(s streamSet, e *editor) *commandline {
	return &commandline{
		streamSet: s,
		editor:    e,
		basic:     &basic{},
	}
}

func (e *commandline) Mode() mode {
	return modeCommandline
}

func (e *commandline) Position() int {
	return e.basic.pos + 1
}

func (e *commandline) Runes() []rune {
	return e.buf
}

func (e *commandline) Message() []rune {
	return append([]rune{':'}, e.basic.buf...)
}

func (e *commandline) Highlight() *screen.Hi {
	return nil
}

func (e *commandline) Run() (end continuity, next modeChanger, err error) {
	r, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return end, next, err
	}
	switch r {
	case CharCtrlM, CharCtrlJ:
	case CharEscape, CharCtrlC:
		next = norm()
		return end, next, err
	case CharBackspace, CharCtrlH:
		if len(e.basic.buf) == 0 {
			next = norm()
			return
		}
		e.basic.delete(e.basic.pos-1, e.basic.pos)
	case CharCtrlB:
		e.basic.move(0)
	case CharCtrlE:
		e.basic.move(len(e.basic.buf))
	case CharCtrlU:
		e.basic.delete(0, e.basic.pos)
	case CharCtrlW:
		// FIXME: It's redundant.
		ed := newEditor()
		ed.pos = e.basic.pos
		ed.buf = e.basic.buf
		pos := ed.pos
		ed.wordBackward()
		e.basic.delete(pos, ed.pos)
		return
	default:
		e.basic.insert([]rune{r}, e.basic.pos)
	}
	if r != CharCtrlM && r != CharCtrlJ {
		return
	}
	next = norm()
	var candidate exCommand
	s := string(e.basic.buf)
	if s == "" {
		return
	}
	for _, cmd := range exCommands {
		if !strings.HasPrefix(cmd.name, s) {
			continue
		}
		if cmd.name == s {
			end = cmd.fn(e, nil)
			return
		}
		if candidate.name == "" {
			candidate = cmd
		}
	}
	if candidate.name != "" {
		end = candidate.fn(e, nil)
		return
	}
	err = fmt.Errorf("not a command: %q", s)
	return
}

func (e *commandline) quit(args []string) continuity {
	return exit
}

func (e *commandline) delete(args []string) (_ continuity) {
	e.editor.delete(0, len(e.editor.buf))
	return
}

func (e *commandline) help(args []string) continuity {
	e.buf = []rune("help")
	e.pos = 4
	return execute
}

func (e *commandline) substitute(args []string) (_ continuity) {
	return
}
