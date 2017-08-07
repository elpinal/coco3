package editor

import (
	"bufio"
	"strings"
	"testing"
)

func TestRecordableRuneReader(t *testing.T) {
	rd := NewReader(bufio.NewReaderSize(strings.NewReader("ABCDEF"), 64))
	r, _, err := rd.ReadRune()
	if err != nil {
		t.Errorf("ReadRune: %v", err)
	}
	if r != 'A' {
		t.Errorf("ReadRune: got %v, want %v", r, 'A')
	}
	rd.Record()
	want := []rune("BCD")
	for i := 0; i < len(want); i++ {
		r, _, err := rd.ReadRune()
		if err != nil {
			t.Errorf("ReadRune: %v", err)
		}
		if r != want[i] {
			t.Errorf("ReadRune: got %q, want %q", r, want[i])
		}
	}
	s := rd.Stop()
	if string(s) != string(want) {
		t.Errorf("record: got %q, want %q", s, want)
	}

	want = []rune("EF")
	for i := 0; i < len(want); i++ {
		r, _, err := rd.ReadRune()
		if err != nil {
			t.Errorf("ReadRune: %v", err)
		}
		if r != want[i] {
			t.Errorf("ReadRune: got %q, want %q", r, want[i])
		}
	}
}
