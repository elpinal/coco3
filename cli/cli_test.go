package cli

import (
	"bytes"
	"testing"

	"github.com/elpinal/coco3/config"
)

func TestFlagC(t *testing.T) {
	var out, err bytes.Buffer
	c := CLI{
		Out: &out,
		Err: &err,
	}
	args := []string{"-c", "echo aaa"}
	code := c.Run(args)
	if code != 0 {
		t.Errorf("Run: got %v, want %v", code, 0)
	}
	if got, want := out.String(), "aaa\n"; got != want {
		t.Errorf("output: got %q, want %q", got, want)
	}
	if e := err.String(); e != "" {
		t.Errorf("error: %v", e)
	}
}

func TestArgs(t *testing.T) {
	var out, err bytes.Buffer
	c := CLI{
		Out: &out,
		Err: &err,
	}
	args := []string{"testdata/basic.coco"}
	code := c.Run(args)
	if code != 0 {
		t.Errorf("Run: got %v, want %v", code, 0)
	}
	if got, want := out.String(), "aaa\nbbb\n"; got != want {
		t.Errorf("output: got %q, want %q", got, want)
	}
	if e := err.String(); e != "" {
		t.Errorf("error: %v", e)
	}
}

func TestExit(t *testing.T) {
	var out, err bytes.Buffer
	c := CLI{
		Out: &out,
		Err: &err,
	}
	args := []string{"-c", "echo aaa; exit 42"}
	code := c.Run(args)
	if code != 42 {
		t.Errorf("Run: got %v, want %v", code, 42)
	}
	if got, want := out.String(), "aaa\n"; got != want {
		t.Errorf("output: got %q, want %q", got, want)
	}
	if e := err.String(); e != "" {
		t.Errorf("error: %v", e)
	}
}

func TestExitInStartUp(t *testing.T) {
	var out, err bytes.Buffer
	c := CLI{
		Out: &out,
		Err: &err,
		Config: config.Config{
			StartUpCommand: []byte("exit 21"),
		},
	}
	args := []string{"-c", "echo aaa; exit 42"}
	code := c.Run(args)
	if code != 21 {
		t.Errorf("Run: got %v, want %v", code, 42)
	}
	if got, want := out.String(), ""; got != want {
		t.Errorf("output: got %q, want %q", got, want)
	}
	if e := err.String(); e != "" {
		t.Errorf("error: %v", e)
	}
}
