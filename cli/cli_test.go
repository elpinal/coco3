package cli

import (
	"bytes"
	"io/ioutil"
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

func TestStartUpCommand(t *testing.T) {
	var out, err bytes.Buffer
	c := CLI{
		Out: &out,
		Err: &err,
		Config: config.Config{
			StartUpCommand: []byte("echo startup..."),
		},
	}
	args := []string{"-c", "echo aaa; echo bbb"}
	code := c.Run(args)
	if code != 0 {
		t.Errorf("Run: got %v, want %v", code, 0)
	}
	if got, want := out.String(), "startup...\naaa\nbbb\n"; got != want {
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

func TestExitInFiles(t *testing.T) {
	var out, err bytes.Buffer
	c := CLI{
		Out: &out,
		Err: &err,
	}
	args := []string{"testdata/exit1.coco", "testdata/exit2.coco"}
	code := c.Run(args)
	if code != 42 {
		t.Errorf("Run: got %v, want %v", code, 42)
	}
	golden, e := ioutil.ReadFile("testdata/exit1.golden")
	if e != nil {
		t.Errorf("reading a golden file: %v", e)
	}
	if got, want := out.String(), string(golden); got != want {
		t.Errorf("output: got %q, want %q", got, want)
	}
	if e := err.String(); e != "" {
		t.Errorf("error: %v", e)
	}
}
