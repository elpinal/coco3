package gate

import (
	"io"

	"github.com/elpinal/coco3/coco4/editor"
	"github.com/elpinal/coco3/coco4/screen/terminal"
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
	b, err := g.e.Read()
	if err != nil {
		return nil, err
	}
	g.history = append(g.history, b)
	return b, nil
}

func (g *gate) Clear() {
	g.e.Clear()
}

func New(in io.Reader, out, err io.Writer) Gate {
	return &gate{
		e: editor.New(terminal.New(out), in, out, err),
	}
}
