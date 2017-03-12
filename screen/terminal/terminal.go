package terminal

import (
	"fmt"
	"io"

	"github.com/mattn/go-runewidth"
)

type Terminal struct {
	w   io.Writer
	msg string
}

func New(w io.Writer) *Terminal {
	return &Terminal{w: w}
}

func (t *Terminal) Refresh(prompt string, s []rune, pos int) {
	io.WriteString(t.w, "\r\033[J")
	io.WriteString(t.w, prompt)
	io.WriteString(t.w, string(s))
	if t.msg != "" {
		io.WriteString(t.w, "\n")
		io.WriteString(t.w, t.msg)
		io.WriteString(t.w, "\033[A")
	}
	fmt.Fprintf(t.w, "\033[%vG", runewidth.StringWidth(prompt)+runesWidth(s[:pos])+1)
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
