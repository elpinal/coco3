package terminal

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/elpinal/coco3/config"
	"github.com/elpinal/coco3/screen"

	"github.com/mattn/go-runewidth"
)

type Terminal struct {
	w              *bufio.Writer
	msg            string
	lastCursorLine int
}

func New(w io.Writer) *Terminal {
	return &Terminal{
		w: bufio.NewWriterSize(w, 32),
	}
}

func getwd() string {
	wd, _ := os.Getwd()
	if home := os.Getenv("HOME"); strings.HasPrefix(wd, home) {
		wd = strings.Replace(wd, home, "~", 1)
	}
	return wd
}

func (t *Terminal) Start(conf *config.Config, inCommandline bool, s []rune, pos int, hi *screen.Hi) {
	t.draw(conf, inCommandline, s, pos, false, hi)
}

func (t *Terminal) Refresh(conf *config.Config, inCommandline bool, s []rune, pos int, hi *screen.Hi) {
	t.draw(conf, inCommandline, s, pos, true, hi)
}

func (t *Terminal) draw(conf *config.Config, inCommandline bool, s []rune, pos int, refresh bool, hi *screen.Hi) {
	prompt := conf.Prompt
	if conf.PromptTmpl != nil {
		var buf bytes.Buffer
		conf.PromptTmpl.Execute(&buf, config.Info{WD: getwd()})
		prompt = buf.String()
	}
	if refresh {
		if t.lastCursorLine > 0 {
			t.w.WriteString("\033[")
			t.w.WriteString(strconv.Itoa(t.lastCursorLine))
			t.w.WriteString("A")
		}
	}
	count := strings.Count(prompt, "\n")
	if inCommandline {
		count++
	}
	t.lastCursorLine = count
	t.w.WriteString("\r\033[J")
	t.w.WriteString(prompt)
	i := strings.LastIndex(prompt, "\n") + 1
	promptWidth := runewidth.StringWidth(prompt[i:])
	if hi == nil {
		t.w.WriteString(string(s))
	} else {
		t.w.WriteString(string(s[:hi.Left]))
		t.w.WriteString("\033[7m")
		t.w.WriteString(string(s[hi.Left:hi.Right]))
		t.w.WriteString("\033[0m")
		t.w.WriteString(string(s[hi.Right:]))
	}
	if t.msg != "" {
		t.w.WriteString("\n\r")
		t.w.WriteString(t.msg)
		if !inCommandline {
			t.w.WriteString("\033[A")
		}
	}
	var drawPos int
	if inCommandline {
		drawPos = runewidth.StringWidth(t.msg[:pos])
	} else {
		drawPos = promptWidth + runesWidth(s[:pos])
	}
	t.w.WriteString("\033[")
	t.w.WriteString(strconv.Itoa(drawPos + 1))
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
