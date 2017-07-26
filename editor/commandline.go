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
	{"quit", (*commandline).quit},
}

type commandline struct {
	streamSet
	*editor

	basic *basic
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

func (e *commandline) Run() (end continuity, next mode, err error) {
	next = modeCommandline
	r, _, err := e.streamSet.in.ReadRune()
	if err != nil {
		return end, next, err
	}
	switch r {
	case CharCtrlM:
	case CharEscape:
		return end, modeNormal, err
	case CharBackspace, CharCtrlH:
		if len(e.basic.buf) == 0 {
			next = modeNormal
			return
		}
		e.basic.delete(e.basic.pos-1, e.basic.pos)
	default:
		e.basic.insert([]rune{r}, e.basic.pos)
	}
	if r != CharCtrlM {
		return
	}
	next = modeNormal
	var candidate exCommand
	s := string(e.basic.buf)
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
