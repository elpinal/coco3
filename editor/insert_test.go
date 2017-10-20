package editor

import (
	"strings"
	"testing"
)

func TestInsertMode(t *testing.T) {
	i := insert{}
	i.init()
	input := "abc"
	i.in = NewReader(strings.NewReader(input))
	for n := range input {
		c, _, err := i.Run()
		if err != nil {
			t.Errorf("Run: %v", err)
		}
		if c != cont {
			t.Errorf("Run: want %v, but got %v", cont, c)
		}
		if got, want := string(i.Runes()), input[:n+1]; got != want {
			t.Errorf("Run: want %v, but got %v", want, got)
		}
	}
}

func TestInputMatches(t *testing.T) {
	i := insert{}
	i.init()
	input := "a 'b"
	i.in = NewReader(strings.NewReader(input))
	tt := []string{"a", "a ", "a ''", "a 'b'"}
	for n := range input {
		c, _, err := i.Run()
		if err != nil {
			t.Errorf("Run: %v", err)
		}
		if c != cont {
			t.Errorf("Run: want %v, but got %v", cont, c)
		}
		if got, want := string(i.Runes()), tt[n]; got != want {
			t.Errorf("Run/%d: want %q, but got %q", n, want, got)
		}
	}
}
