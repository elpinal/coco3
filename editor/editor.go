package editor

import (
	"bufio"
	"context"
	"io"

	"github.com/elpinal/coco3/config"
	"github.com/elpinal/coco3/screen"
)

type Editor interface {
	Read() ([]rune, bool, error)
	Clear()
	SetHistory([][]rune)
}

func New(s screen.Screen, conf *config.Config, in io.Reader, out, err io.Writer) Editor {
	return NewContext(context.Background(), s, conf, in, out, err)
}

func NewContext(ctx context.Context, s screen.Screen, conf *config.Config, in io.Reader, out, err io.Writer) Editor {
	var rd io.RuneReader
	if x, ok := in.(io.RuneReader); ok {
		rd = x
	} else {
		rd = bufio.NewReaderSize(in, 64)
	}
	r := NewReaderContext(ctx, rd)
	return &balancer{
		streamSet: streamSet{
			in:  r,
			out: out,
			err: err,
		},
		editor: newEditor(),
		s:      s,
		conf:   conf,
	}
}

type streamSet struct {
	in  *RuneAddReader
	out io.Writer
	err io.Writer
}

type balancer struct {
	streamSet
	*editor
	s    screen.Screen
	conf *config.Config
}

func (b *balancer) Read() ([]rune, bool, error) {
	b.s.SetLastLine("-- INSERT --")
	b.s.Start(b.conf, false, nil, 0)
	prev := modeInsert
	m := b.enter(prev)
	for {
		end, next, err := m.Run()
		if err != nil {
			return nil, false, err
		}
		if end == exit {
			return nil, true, nil
		}
		if end == execute {
			b.s.SetLastLine("")
			b.s.Refresh(b.conf, m.Mode() == modeCommandline, m.Runes(), m.Position())
			return m.Runes(), false, nil
		}
		if prev != next {
			m = b.enter(next)
		}
		prev = next
		msg := string(m.Message())
		b.s.SetLastLine(msg)
		b.s.Refresh(b.conf, m.Mode() == modeCommandline, m.Runes(), m.Position())
	}
}

func (b *balancer) enter(m mode) moder {
	switch m {
	case modeInsert:
		return &insert{
			streamSet: b.streamSet,
			editor:    b.editor,
			s:         b.s,
			conf:      b.conf,
		}
	case modeNormal:
		return newNormal(
			b.streamSet,
			b.editor,
		)
	case modeVisual:
		return newVisual(
			b.streamSet,
			b.editor,
		)
	case modeReplace:
		buf := b.buf
		b.buf = nil
		return &insert{
			streamSet:   b.streamSet,
			editor:      b.editor,
			replaceMode: true,
			replacedBuf: buf,
		}
	case modeCommandline:
		return &commandline{
			streamSet: b.streamSet,
			editor:    b.editor,
			basic:     &basic{},
		}
	}
	return &insert{streamSet: b.streamSet}
}

func (b *balancer) Clear() {
	b.buf = b.buf[:0]
	b.pos = 0
}

func (b *balancer) SetHistory(history [][]rune) {
	// copy history
	b.history = make([][]rune, len(history))
	for i, h := range history {
		b.history[i] = make([]rune, len(h))
		copy(b.history[i], h)
	}
	b.age = len(history)
}

const (
	mchar = iota
	mline
)

type continuity int

const (
	cont continuity = iota
	execute
	exit
)
