package editor

import (
	"bufio"
	"strings"
	"testing"
)

func TestRuneAddReader(t *testing.T) {
	rd := NewReader(bufio.NewReaderSize(strings.NewReader("ABCDEF"), 64))
	r, _, err := rd.ReadRune()
	if err != nil {
		t.Errorf("ReadRune: %v", err)
	}
	if r != 'A' {
		t.Errorf("ReadRune: got %v, want %v", r, 'A')
	}
	rd.Add([]rune("ab"))
	want := []rune("abBCDEF")
	for i := 0; i < 7; i++ {
		r, _, err = rd.ReadRune()
		if err != nil {
			t.Errorf("ReadRune: %v", err)
		}
		if r != want[i] {
			t.Errorf("ReadRune: got %v, want %v", string(r), string(want[i]))
		}
	}
}

func TestRuneAddReaderMultiple(t *testing.T) {
	rd := NewReader(bufio.NewReaderSize(strings.NewReader(strings.Repeat("a", 128)), 64))
	r, _, err := rd.ReadRune()
	if err != nil {
		t.Errorf("ReadRune: %v", err)
	}
	if r != 'a' {
		t.Errorf("ReadRune: got %v, want %v", r, 'a')
	}
	for range [...]int{127: 0} {
		rd.Add([]rune("dl"))
	}
	want := []rune(strings.Repeat("dl", 128) + strings.Repeat("a", 127))
	for i := 0; i < len(want); i++ {
		r, _, err = rd.ReadRune()
		if err != nil {
			t.Errorf("ReadRune: %v", err)
		}
		if r != want[i] {
			t.Errorf("ReadRune: got %v, want %v", string(r), string(want[i]))
		}
	}
}
