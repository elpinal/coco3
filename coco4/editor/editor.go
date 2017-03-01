package editor

import (
	"bufio"
	"io"

	"github.com/elpinal/coco3/coco4/screen"
)

type Editor interface {
	Read() ([]rune, error)
	Clear()
}

func New(s screen.Screen, in io.Reader, out, err io.Writer) Editor {
	var rd io.RuneReader
	if x, ok := in.(io.RuneReader); ok {
		rd = x
	} else {
		rd = bufio.NewReaderSize(in, 64)
	}
	return &balancer{
		streamSet: streamSet{
			in:  rd,
			out: out,
			err: err,
		},
		editor: &editor{},
		s:      s,
	}
}

type streamSet struct {
	in  io.RuneReader
	out io.Writer
	err io.Writer
}

type balancer struct {
	streamSet
	*editor
	s screen.Screen
}

func (b *balancer) Read() ([]rune, error) {
	prev := modeInsert
	m := b.enter(prev)
	for {
		end, next, err := m.Run()
		if err != nil {
			return nil, err
		}
		b.s.Refresh(string(m.Runes()), m.Position())
		if end {
			return m.Runes(), nil
		}
		if prev != next {
			m = b.enter(next)
		}
		prev = next
	}
}

func (b *balancer) enter(m mode) moder {
	switch m {
	case modeInsert:
		return &insert{
			streamSet: streamSet{
				in:  b.in,
				out: b.out,
				err: b.err,
			},
			editor: b.editor,
		}
	case modeNormal:
		return &normal{
			streamSet: streamSet{
				in:  b.in,
				out: b.out,
				err: b.err,
			},
			editor: b.editor,
		}
	}
	return &insert{streamSet: b.streamSet}
}

func (b *balancer) Clear() {
	b.buf = b.buf[:0]
	b.pos = 0
}

const (
	mchar = iota
	mline
)
