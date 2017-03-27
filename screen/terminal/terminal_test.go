package terminal

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestTerminal(t *testing.T) {
	var buf bytes.Buffer
	term := New(&buf)
	term.SetLastLine("-- last line --")
	term.Refresh("prompt ", []rune("aaa"), 2)
	if got := buf.String(); strings.Index(got, "prompt aaa") < 0 {
		t.Errorf("got %q, but should include %q", got, "prompt aaa")
	}
}

func BenchmarkTerminal(b *testing.B) {
	term := New(ioutil.Discard)
	term.SetLastLine("-- last line --")
	for i := 0; i < b.N; i++ {
		term.Refresh("prompt ", []rune("aaa"), 2)
	}
}
