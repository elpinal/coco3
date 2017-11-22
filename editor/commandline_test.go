package editor

import (
	"bytes"
	"testing"
)

func TestCommandline(t *testing.T) {
	command := append([]byte("quit"), CharCtrlM)
	ed := newCommandline(
		streamSet{in: NewReader(bytes.NewReader(command))},
		newEditor(),
	)
	var (
		end continuity
		err error
	)
	for range command {
		end, _, err = ed.Run()
		if err != nil {
			t.Errorf("commandline: %v", err)
		}
	}
	if end != exit {
		t.Errorf("commandline (%q): want %v, but got %v", command, exit, end)
	}
}

func TestSubstituteNoArgs(t *testing.T) {
	command := append([]byte("substitute"), CharCtrlM)
	ed := newCommandline(
		streamSet{in: NewReader(bytes.NewReader(command))},
		newEditor(),
	)
	ed.buf = []rune("a")
	for range command {
		_, _, err := ed.Run()
		if err != nil {
			t.Errorf("commandline: %v", err)
		}
	}
	got := string(ed.buf)
	want := "a"
	if got != want {
		t.Errorf("commandline (%q): got %v, want %v", command, got, want)
	}
}

func TestSubstitute(t *testing.T) {
	command := append([]byte("substitute a b"), CharCtrlM)
	ed := newCommandline(
		streamSet{in: NewReader(bytes.NewReader(command))},
		newEditor(),
	)
	ed.buf = []rune("a")
	for range command {
		_, _, err := ed.Run()
		if err != nil {
			t.Errorf("commandline: %v", err)
		}
	}
	got := string(ed.buf)
	want := "b"
	if got != want {
		t.Errorf("commandline (%q): got %v, want %v", command, got, want)
	}
}

func TestEmpty(t *testing.T) {
	tests := []struct {
		command []byte
	}{
		{command: append([]byte(""), CharCtrlM)},
		{command: append([]byte(" "), CharCtrlM)},
		{command: append([]byte("  "), CharCtrlM)},
	}
	for _, test := range tests {
		ed := newCommandline(
			streamSet{in: NewReader(bytes.NewReader(test.command))},
			newEditor(),
		)
		for range test.command {
			_, _, err := ed.Run()
			if err != nil {
				t.Errorf("commandline: %v", err)
			}
		}
	}
}

func TestSubstituteSpaces(t *testing.T) {
	command := append([]byte(`substitute " a " "_v_"`), CharCtrlM)
	ed := newCommandline(
		streamSet{in: NewReader(bytes.NewReader(command))},
		newEditor(),
	)
	ed.buf = []rune(" a  a ")
	for range command {
		_, _, err := ed.Run()
		if err != nil {
			t.Errorf("commandline: %v", err)
		}
	}
	got := string(ed.buf)
	want := "_v__v_"
	if got != want {
		t.Errorf("commandline (%q): got %v, want %v", command, got, want)
	}
}
