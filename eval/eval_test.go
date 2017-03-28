package eval

import (
	"bytes"
	"syscall"
	"testing"
	"time"
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

type myReader struct{}

func (r *myReader) Read(_ []byte) (int, error) {
	return 0, nil
}

func TestInterrupt(t *testing.T) {
	in := myReader{}
	var out, err bytes.Buffer
	e := New(&in, &out, &err)
	done := make(chan struct{})
	go func() {
		if err := e.execCmd("cat", nil); err != nil && err == ErrInterrupted {
		} else if err == nil {
			t.Error("cat: error should not be nil because the command must have been interrupted")
		} else {
			t.Errorf("cat: %v", err)
		}
		done <- struct{}{}
	}()
	time.Sleep(50 * time.Millisecond)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	<-done
}
