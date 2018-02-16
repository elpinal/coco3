package editor

import (
	"strings"
	"testing"
)

func BenchmarkNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		command := strings.Repeat("hjkl", 100)
		ed := newNormal(
			streamSet{
				in: NewReader(strings.NewReader(command)),
			},
			newEditor(),
		)
		for range command {
			_, _, err := ed.Run()
			if err != nil {
				b.Errorf("normal: %v", err)
			}
		}
	}
}

func TestIndexNumber(t *testing.T) {
	norm := newNormal(streamSet{in: NewReader(nil)}, newEditor())
	if got := norm.indexNumber(); got >= 0 {
		t.Errorf("got %d, but want -1", got)
	}

	norm = newNormal(streamSet{in: NewReader(nil)}, newEditor())
	norm.buf = []rune("4") // TODO: this way of creating "normal" is not sophisticated.
	if got := norm.indexNumber(); got != 0 {
		t.Errorf("got %d, but want 0", got)
	}

	norm = newNormal(streamSet{in: NewReader(nil)}, newEditor())
	norm.buf = []rune("-a4")
	if got, want := norm.indexNumber(), 2; got != want {
		t.Errorf("got %d, but want %d", got, want)
	}
}

func TestUpdateNumber(t *testing.T) {
	norm := newNormal(streamSet{in: NewReader(nil)}, newEditor())
	norm.updateNumber(func(n int) int {
		return n
	})

	norm = newNormal(streamSet{in: NewReader(nil)}, newEditor())
	norm.buf = []rune("4") // TODO: this way of creating "normal" is not sophisticated.
	norm.updateNumber(func(n int) int {
		if n != 4 {
			t.Fatalf("got %d, but want 4", n)
		}
		return n + 1
	})
	if got := string(norm.buf); got != "5" {
		t.Errorf("got %q, but want %q", got, "5")
	}

	norm = newNormal(streamSet{in: NewReader(nil)}, newEditor())
	norm.buf = []rune("-a4")
	norm.updateNumber(func(n int) int {
		if n != 4 {
			t.Fatalf("got %d, but want 4", n)
		}
		return n + 1
	})
	if got := string(norm.buf); got != "-a5" {
		t.Errorf("got %q, but want %q", got, "-a5")
	}
}
