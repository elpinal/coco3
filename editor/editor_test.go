package editor

import (
	"bytes"
	"strings"
	"testing"

	"github.com/elpinal/coco3/config"
)

type testScreen struct {
}

func (ts *testScreen) Refresh(prompt string, s []rune, pos int) {
}

func (ts *testScreen) SetLastLine(msg string) {
}

func TestEditor(t *testing.T) {
	inBuf := strings.NewReader("aaa" + string(CharCtrlM))
	var outBuf, errBuf bytes.Buffer
	e := New(&testScreen{}, &config.Config{}, inBuf, &outBuf, &errBuf)
	s, err := e.Read()
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
