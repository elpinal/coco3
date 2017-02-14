package main

import "testing"

func TestUndoTree(t *testing.T) {
	s1 := "The first history"
	s2 := "The second history"
	s3 := "The third history"
	u := newUndoTree()

	u.add([]rune(s1))
	if got, want := string(u.undo()), ""; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	u.add([]rune(s2))
	u.add([]rune(s3))
	if got, want := string(u.undo()), s2; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	u.add([]rune(s2))
	if got, want := string(u.undo()), s2; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(u.redo()), s2; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(u.later()), ""; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(u.earlier()), s3; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(u.earlier()), s2; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(u.earlier()), s1; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := string(u.earlier()), ""; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
