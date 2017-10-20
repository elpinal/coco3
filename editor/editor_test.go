package editor

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/elpinal/coco3/config"
	"github.com/elpinal/coco3/screen"
)

func TestEditor(t *testing.T) {
	inBuf := strings.NewReader("aaa" + string(CharCtrlM))
	var outBuf, errBuf bytes.Buffer
	e := New(&screen.TestScreen{}, &config.Config{}, inBuf, &outBuf, &errBuf)
	s, _, err := e.Read()
	if err != nil {
		t.Error(err)
	}
	if want := "aaa"; string(s) != want {
		t.Errorf("got %q, want %q", string(s), want)
	}
	if got := outBuf.String(); got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
	if got := errBuf.String(); got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
	e.Clear()
}

func TestNormal(t *testing.T) {
	inBuf := strings.NewReader("aaa" + string([]rune{
		CharEscape,
		'3', 'h',
		'2', 'x',
		'i',
		'A',
		CharEscape,
		'y', 'y',
		'2', 'p',
		'i',
		CharCtrlM,
	}))
	var outBuf, errBuf bytes.Buffer
	e := New(&screen.TestScreen{}, &config.Config{}, inBuf, &outBuf, &errBuf)
	s, _, err := e.Read()
	if err != nil {
		t.Error(err)
	}
	if want := "AAaAaa"; string(s) != want {
		t.Errorf("got %q, want %q", string(s), want)
	}
	if got := outBuf.String(); got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
	if got := errBuf.String(); got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
	e.Clear()
}

func TestWithVisual(t *testing.T) {
	input := strings.NewReader("aaa" + string([]rune{
		CharEscape,
		'0',
		'l',
		'v',
		'A',
		'X',
		CharEscape,
		'V',
		's', '`',
		'i',
		CharCtrlM,
	}))
	var outBuf, errBuf bytes.Buffer
	e := New(&screen.TestScreen{}, &config.Config{}, input, &outBuf, &errBuf)
	s, _, err := e.Read()
	if err != nil {
		t.Error(err)
	}
	if want := "`aaXa`"; string(s) != want {
		t.Errorf("got %q, want %q", string(s), want)
	}
	if got := outBuf.String(); got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
	if got := errBuf.String(); got != "" {
		t.Errorf("got %q, want %q", got, "")
	}
	e.Clear()
}

func TestYankPaste(t *testing.T) {
	e := New(&screen.TestScreen{}, &config.Config{}, strings.NewReader("123"+string([]rune{CharEscape, 'h', 'x', 'p', 'i', CharCtrlM})), ioutil.Discard, ioutil.Discard)
	s, _, err := e.Read()
	if err != nil {
		t.Error(err)
	}
	if want := "132"; string(s) != want {
		t.Errorf("got %q, want %q", string(s), want)
	}
}
