package gate

import (
	"context"
	"io"

	"github.com/elpinal/coco3/config"
	"github.com/elpinal/coco3/editor"
	"github.com/elpinal/coco3/screen/terminal"
)

type Gate interface {
	Read() ([]rune, bool, error)
	Clear()
}

type gate struct {
	e editor.Editor

	history [][]rune
}

func (g *gate) Read() ([]rune, bool, error) {
	g.e.SetHistory(g.history)
	b, end, err := g.e.Read()
	if err != nil {
		return nil, false, err
	}
	if end {
		return nil, true, nil
	}
	if len(b) != 0 && (len(g.history) == 0 || string(g.history[len(g.history)-1]) != string(b)) {
		g.history = append(g.history, b)
	}
	return b, false, nil
}

func (g *gate) Clear() {
	g.e.Clear()
}

func New(conf *config.Config, in io.Reader, out, err io.Writer, history [][]rune) Gate {
	return NewContext(context.Background(), conf, in, out, err, history)
}

func NewContext(ctx context.Context, conf *config.Config, in io.Reader, out, err io.Writer, history [][]rune) Gate {
	return &gate{
		e:       editor.NewContext(ctx, terminal.New(out), conf, in, out, err),
		history: history,
	}
}
