package terminal

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/elpinal/coco3/config"
)

func TestTerminal(t *testing.T) {
	var buf bytes.Buffer
	term := New(&buf)
	term.SetLastLine("-- last line --")
	term.Start(&config.Config{Prompt: "prompt "}, false, []rune("aaa"), 2, nil)
	if got := buf.String(); !strings.Contains(got, "prompt aaa") {
		t.Errorf("got %q, but should include %q", got, "prompt aaa")
	}
}

func BenchmarkTerminal(b *testing.B) {
	term := New(ioutil.Discard)
	term.SetLastLine("-- last line --")
	for i := 0; i < b.N; i++ {
		term.Refresh(&config.Config{Prompt: "prompt "}, false, []rune("aaa"), 2, nil)
	}
}
