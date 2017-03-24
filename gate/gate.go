package gate

import (
	"io"

	"github.com/elpinal/coco3/config"
	"github.com/elpinal/coco3/editor"
	"github.com/elpinal/coco3/screen/terminal"
)

type Gate interface {
	Read() ([]rune, error)
	Clear()
}

type gate struct {
	e editor.Editor

	history [][]rune
}

func (g *gate) Read() ([]rune, error) {
	g.e.SetHistory(g.history)
	b, err := g.e.Read()
	if err != nil {
		return nil, err
	}
	if len(g.history) == 0 || string(g.history[len(g.history)-1]) != string(b) {
		g.history = append(g.history, b)
	}
	return b, nil
}

func (g *gate) Clear() {
	g.e.Clear()
}

func New(conf *config.Config, in io.Reader, out, err io.Writer) Gate {
	return &gate{
		e: editor.New(terminal.New(out), conf, in, out, err),
	}
}
