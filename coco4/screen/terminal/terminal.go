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

func (t *Terminal) Refresh(s string, pos int) {
	io.WriteString(t.w, "\r\033[K")
	io.WriteString(t.w, s)
	fmt.Fprintf(t.w, "\033[%vG", runewidth.StringWidth(s[:pos])+1)

}
