package terminal

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/elpinal/coco3/config"

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

func getwd() string {
	wd, _ := os.Getwd()
	if home := os.Getenv("HOME"); strings.HasPrefix(wd, home) {
		wd = strings.Replace(wd, home, "~", 1)
	}
	return wd
}

func (t *Terminal) Start(conf *config.Config, s []rune, pos int) {
	var promptWidth int
	if conf.PromptTmpl != nil {
		var buf bytes.Buffer
		conf.PromptTmpl.Execute(&buf, config.Info{WD: getwd()})
		prompt := buf.String()
		t.w.WriteString("\r\033[J")
		t.w.WriteString(prompt)
		i := strings.LastIndex(prompt, "\n") + 1
		lastLine := prompt[i:]
		promptWidth = runewidth.StringWidth(lastLine)
	} else {
		t.w.WriteString("\r\033[J")
		t.w.WriteString(conf.Prompt)
		promptWidth = runewidth.StringWidth(conf.Prompt)
	}
	t.w.WriteString(string(s))
	if t.msg != "" {
		t.w.WriteString("\n")
		t.w.WriteString(t.msg)
		t.w.WriteString("\033[A")
	}
	t.w.WriteString("\033[")
	t.w.WriteString(strconv.Itoa(promptWidth + runesWidth(s[:pos]) + 1))
	t.w.WriteString("G")
	t.w.Flush()
}

func (t *Terminal) Refresh(conf *config.Config, s []rune, pos int) {
	var promptWidth int
	if conf.PromptTmpl != nil {
		var buf bytes.Buffer
		conf.PromptTmpl.Execute(&buf, config.Info{WD: getwd()})
		prompt := buf.String()
		count := strings.Count(prompt, "\n")
		if count > 0 {
			t.w.WriteString("\033[")
			t.w.WriteString(strconv.Itoa(count))
			t.w.WriteString("A")
		}
		t.w.WriteString("\r\033[J")
		t.w.WriteString(prompt)
		i := strings.LastIndex(prompt, "\n") + 1
		lastLine := prompt[i:]
		promptWidth = runewidth.StringWidth(lastLine)
	} else {
		t.w.WriteString("\r\033[J")
		t.w.WriteString(conf.Prompt)
		promptWidth = runewidth.StringWidth(conf.Prompt)
	}
	t.w.WriteString(string(s))
	if t.msg != "" {
		t.w.WriteString("\n")
		t.w.WriteString(t.msg)
		t.w.WriteString("\033[A")
	}
	t.w.WriteString("\033[")
	t.w.WriteString(strconv.Itoa(promptWidth + runesWidth(s[:pos]) + 1))
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
