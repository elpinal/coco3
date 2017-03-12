package main

import (
	"bytes"
	"testing"
)

func TestFlagC(t *testing.T) {
	var in, out, err bytes.Buffer
	c := cli{
		in:  &in,
		out: &out,
		err: &err,
	}
	args := []string{"-c", "echo aaa"}
	code := c.run(args)
	if code != 0 {
		t.Errorf("run: got %v, want %v", code, 0)
	}
	if got, want := out.String(), "aaa\n"; got != want {
		t.Errorf("output: got %v, want %v", got, want)
	}
}

func TestArgs(t *testing.T) {
	var in, out, err bytes.Buffer
	c := cli{
		in:  &in,
		out: &out,
		err: &err,
	}
	args := []string{"testdata/basic.coco"}
	code := c.run(args)
	if code != 0 {
		t.Errorf("run: got %v, want %v", code, 0)
	}
	if got, want := out.String(), "aaa\nbbb\n"; got != want {
		t.Errorf("output: got %v, want %v", got, want)
	}
}
