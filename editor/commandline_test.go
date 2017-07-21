package editor

import (
	"bytes"
	"context"
	"testing"
)

func TestCommandline(t *testing.T) {
	command := append([]byte("quit"), CharCtrlM)
	ed := commandline{
		streamSet: streamSet{in: NewReaderContext(context.TODO(), bytes.NewReader(command))},
		editor:    newEditor(),
		basic:     &basic{},
	}
	var (
		end  continuity
		next mode
		err  error
	)
	for range command {
		end, next, err = ed.Run()
		if err != nil {
			t.Errorf("commandline: %v", err)
		}
	}
	if end != exit {
		t.Errorf("commandline (%q): want %v, but got %v", command, exit, end)
	}
	if next != modeNormal {
		t.Errorf("commandline (%q): next mode: want %v, but got %v", command, modeNormal, next)
	}
}
