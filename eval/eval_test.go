package eval

import (
	"bytes"
	"testing"
)

func TestExecCmd(t *testing.T) {
	var in, out, err bytes.Buffer
	e := New(&in, &out, &err)
	if err := e.execCmd("echo", []string{"aaa"}); err != nil {
		t.Errorf("execute command: %v", err)
	}
	if got, want := out.String(), "aaa\n"; got != want {
		t.Errorf("output: got %q, want %q", got, want)
	}
	if got := err.String(); got != "" {
		t.Errorf("error should be blank; got %q", got)
	}
}

func TestExecPipe(t *testing.T) {
	var in, out, err bytes.Buffer
	e := New(&in, &out, &err)
	if err := e.execPipe([][]string{
		[]string{
			"echo",
			"aaa",
		},
		[]string{
			"tr",
			"a",
			"A",
		},
	}); err != nil {
		t.Errorf("execute pipe: %v", err)
	}
	if got, want := out.String(), "AAA\n"; got != want {
		t.Errorf("output: got %q, want %q", got, want)
	}
	if got := err.String(); got != "" {
		t.Errorf("error should be blank; got %q", got)
	}
}
