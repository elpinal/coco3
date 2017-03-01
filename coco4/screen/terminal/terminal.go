package terminal

import (
	"fmt"
	"io"

	"github.com/mattn/go-runewidth"
)

type Terminal struct {
	w io.Writer
}

func New(w io.Writer) *Terminal {
	return &Terminal{w: w}
}

func (t *Terminal) Refresh(prompt string, s []rune, pos int) {
	io.WriteString(t.w, "\r\033[K")
	io.WriteString(t.w, prompt)
	io.WriteString(t.w, string(s))
	fmt.Fprintf(t.w, "\033[%vG", runewidth.StringWidth(prompt)+runesWidth(s[:pos])+1)
}

func runesWidth(s []rune) (width int) {
	for _, r := range s {
		width += runewidth.RuneWidth(r)
	}
	return width
}
