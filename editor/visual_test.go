package editor

import (
	"strings"
	"testing"
)

func TestVisual(t *testing.T) {
	input := "hlwb"
	v := newVisual(streamSet{in: NewReader(strings.NewReader(input))}, &editor{basic: basic{buf: []rune("aaa bbb ccc"), pos: 5}})
	for i := range input {
		end, next, err := v.Run()
		if err != nil {
			t.Errorf("Run (%d): %v", i, err)
		}
		if next != nil {
			t.Errorf("Run (%d): got %v, want %v", i, next, nil)
		}
		if end != cont {
			t.Errorf("Run (%d): got %v, want %v", i, end, cont)
		}
	}
	if v.pos != 4 {
		t.Errorf("Run: v.pos should be %d, but got %d", 4, v.pos)
	}
}
