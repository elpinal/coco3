package editor

import (
	"strings"

	parser "github.com/elpinal/coco3/editor/commandline"
	"github.com/elpinal/coco3/screen"
)

type exCommand struct {
	name string
	fn   func(*commandline, []parser.Token) continuity
}

// exComands represents a table of Ex commands and corresponding functions.
// The order is important. Preceding commands have higher precedence.
var exCommands = []exCommand{
	{"help", (*commandline).help},
	{"delete", (*commandline).delete},
	{"quit", (*commandline).quit},
	{"substitute", (*commandline).substitute},
}

type commandline struct {
	streamSet
	*editor

	basic *basic

	// FIXME: currently live a moment.
	history [][]rune
	age     int
}

func emptyCommandline() *commandline {
	return &commandline{
		editor: newEditor(),
		basic: &basic{},
	}
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
	r, _, err := e.in.ReadRune()
	if err != nil {
		return end, next, err
	}
	switch r {
	case CharCtrlM, CharCtrlJ:
		end, err = e.execute()
		next = norm()

	case CharEscape, CharCtrlC:
		next = norm()

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
	case CharCtrlN:
		e.historyForward()
	case CharCtrlP:
		e.historyBack()
	case CharCtrlU:
		e.basic.delete(0, e.basic.pos)
	case CharCtrlW:
		pos := e.basic.pos
		e.basic.wordBackward()
		e.basic.delete(pos, e.basic.pos)
	default:
		e.basic.insert([]rune{r}, e.basic.pos)
	}
	return
}

type ErrNoExCommand struct {
	Name string
}

func (e *ErrNoExCommand) Error() string {
	return "no such Ex command: " + e.Name
}

func (e *commandline) execute() (end continuity, err error) {
	command, err := parser.ParseT(string(e.basic.buf))
	if err != nil {
		return cont, err
	}
	if command == nil {
		return
	}
	name := string(command.Name)
	defer func() {
		e.history = append(e.history, e.basic.buf)
	}()
	var candidate exCommand
	for _, cmd := range exCommands {
		if !strings.HasPrefix(cmd.name, name) {
			continue
		}
		if cmd.name == name {
			end = cmd.fn(e, command.Args)
			return
		}
		if candidate.name == "" {
			candidate = cmd
		}
	}
	if candidate.name != "" {
		end = candidate.fn(e, command.Args)
		return
	}
	err = &ErrNoExCommand{Name: name}
	return
}

func (e *commandline) historyBack() {
	l := len(e.history)
	if e.age == l {
		return
	}
	e.age++
	e.basic.buf = e.history[l-e.age]
}

func (e *commandline) historyForward() {
	l := len(e.history)
	if e.age == 0 {
		return
	}
	e.age--
	e.basic.buf = e.history[l-e.age]
}

func (e *commandline) quit(_ []parser.Token) continuity {
	return exit
}

func (e *commandline) delete(_ []parser.Token) (_ continuity) {
	e.editor.delete(0, len(e.editor.buf))
	return
}

func (e *commandline) help(_ []parser.Token) continuity {
	e.buf = []rune("help")
	e.pos = 4
	return execute
}

func (e *commandline) substitute(args []parser.Token) (_ continuity) {
	if len(args) != 2 {
		return
	}
	pat := toString(args[0])
	s0 := toString(args[1])
	s := strings.Replace(string(e.buf), string(pat), string(s0), -1)
	e.buf = []rune(s)
	return
}

func toString(t parser.Token) string {
	switch t.Type {
	case parser.TokenIdent:
	case parser.TokenString:
		return unquote(string(t.Value))
	}
	return string(t.Value)
}

func unquote(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, `"`), `"`)
}
