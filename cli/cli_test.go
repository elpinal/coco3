package cli

import (
	"bytes"
	"testing"
)

func TestFlagC(t *testing.T) {
	var in, out, err bytes.Buffer
	c := CLI{
		In:  &in,
		Out: &out,
		Err: &err,
	}
	args := []string{"-c", "echo aaa"}
	code := c.Run(args)
	if code != 0 {
		t.Errorf("Run: got %v, want %v", code, 0)
	}
	if got, want := out.String(), "aaa\n"; got != want {
		t.Errorf("output: got %v, want %v", got, want)
	}
	if e := err.String(); e != "" {
		t.Errorf("error: %v", e)
	}
}

func TestArgs(t *testing.T) {
	var in, out, err bytes.Buffer
	c := CLI{
		In:  &in,
		Out: &out,
		Err: &err,
	}
	args := []string{"testdata/basic.coco"}
	code := c.Run(args)
	if code != 0 {
		t.Errorf("Run: got %v, want %v", code, 0)
	}
	if got, want := out.String(), "aaa\nbbb\n"; got != want {
		t.Errorf("output: got %v, want %v", got, want)
	}
	if e := err.String(); e != "" {
		t.Errorf("error: %v", e)
	}
}
