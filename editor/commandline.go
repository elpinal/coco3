package editor

import (
	"fmt"
	"strings"
)

type exCommand struct {
	name string
	fn   func(*commandline, []string) mode
}

// exComands represents a table of Ex commands and corresponding functions.
// The order is important. Precede commands have higher precedence.
var exCommands = []exCommand{
	{"quit", (*commandline).quit},
}

type commandline struct {
	streamSet
	*editor
}

func (e *commandline) Mode() mode {
	return modeCommandline
}

func (e *commandline) Position() int {
	return e.pos
}

func (e *commandline) Runes() []rune {
	return e.buf
}

func (e *commandline) Run() (end bool, next mode, err error) {
	next = modeNormal
	var r rune
	rs := make([]rune, 0, 4)
	for {
		r, _, err = e.streamSet.in.ReadRune()
		if err != nil {
			return
		}
		if r == CharCtrlM {
			break
		}
		rs = append(rs, r)
	}
	s := string(rs)
	if s == "" {
		return
	}
	var candidate exCommand
	for _, cmd := range exCommands {
		if !strings.HasPrefix(cmd.name, s) {
			continue
		}
		if cmd.name == s {
			next = cmd.fn(e, nil)
			return
		}
		if candidate.name == "" {
			candidate = cmd
		}
	}
	if candidate.name != "" {
		next = candidate.fn(e, nil)
		return
	}
	return false, modeCommandline, fmt.Errorf("not a command: %q", s)
}

func (e *commandline) quit(args []string) mode {
	return modeNormal
}
