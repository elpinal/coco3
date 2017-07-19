package editor

import (
	"fmt"
	"strings"
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

func (e *commandline) Run() (end continuity, next mode, err error) {
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