package eval

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"
	"time"
)

func TestExecCmd(t *testing.T) {
	var out, err bytes.Buffer
	e := New(nil, &out, &err)
	if err := e.execCmd(context.Background(), "echo", []string{"aaa"}); err != nil {
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
	var out, err bytes.Buffer
	e := New(nil, &out, &err)
	if err := e.execPipe(context.Background(), [][]string{
		{
			"echo",
			"aaa",
		},
		{
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

type slowWriter struct{}

func (w *slowWriter) Write(p []byte) (int, error) {
	time.Sleep(10 * time.Millisecond)
	return 1, nil
}

func TestKillBuiltin(t *testing.T) {
	e := New(nil, &slowWriter{}, ioutil.Discard)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	if err := e.execCmd(ctx, "echo", []string{100: ""}); err != nil {
		t.Errorf("echo: %v", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("echo: should be killed by 1 second, but elapsed time is %v", elapsed)
	}
}
