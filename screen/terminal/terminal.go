package terminal

import (
	"bytes"
	"io"
	"strconv"

	"github.com/mattn/go-runewidth"
)

type Terminal struct {
	w   io.Writer
	buf *bytes.Buffer
	msg string
}

func New(w io.Writer) *Terminal {
	return &Terminal{
		w:   w,
		buf: bytes.NewBuffer(make([]byte, 0, 32)),
	}
}

func (t *Terminal) Refresh(prompt string, s []rune, pos int) {
	t.buf.WriteString("\r\033[J")
	t.buf.WriteString(prompt)
	t.buf.WriteString(string(s))
	if t.msg != "" {
		t.buf.WriteString("\n")
		t.buf.WriteString(t.msg)
		t.buf.WriteString("\033[A")
	}
	t.buf.WriteString("\033[")
	t.buf.WriteString(strconv.Itoa(runewidth.StringWidth(prompt) + runesWidth(s[:pos]) + 1))
	t.buf.WriteString("G")
	t.buf.WriteTo(t.w)
}

func runesWidth(s []rune) (width int) {
	for _, r := range s {
		width += runewidth.RuneWidth(r)
	}
	return width
}

func (t *Terminal) SetLastLine(msg string) {
	t.msg = msg
}
