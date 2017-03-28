package terminal

import (
	"bufio"
	"io"
	"strconv"

	"github.com/mattn/go-runewidth"
)

type Terminal struct {
	w   *bufio.Writer
	msg string
}

func New(w io.Writer) *Terminal {
	return &Terminal{
		w: bufio.NewWriterSize(w, 32),
	}
}

func (t *Terminal) Refresh(prompt string, s []rune, pos int) {
	t.w.WriteString("\r\033[J")
	t.w.WriteString(prompt)
	t.w.WriteString(string(s))
	if t.msg != "" {
		t.w.WriteString("\n")
		t.w.WriteString(t.msg)
		t.w.WriteString("\033[A")
	}
	t.w.WriteString("\033[")
	t.w.WriteString(strconv.Itoa(runewidth.StringWidth(prompt) + runesWidth(s[:pos]) + 1))
	t.w.WriteString("G")
	t.w.Flush()
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
