package main

import (
	"reflect"
	"testing"
)

func TestMove(t *testing.T) {
	s := []rune("aaaaa bbbbb ccccc")
	tests := []struct {
		input int
		want  int
	}{
		{0, 0},
		{1, 1},
		{10, 10},
		{len(s), len(s)},
		{-1, 0},
		{len(s) + 1, len(s)},
	}
	for _, test := range tests {
		e := &basicEditor{buf: s}
		e.move(test.input)
		if e.pos != test.want {
			t.Errorf("move(%v): got %v, want %v", test.input, e.pos, test.want)
		}
	}
}

func TestInsert(t *testing.T) {
	tests := []struct {
		initial basicEditor
		input   []rune
		at      int
		want    basicEditor
	}{
		{
			initial: basicEditor{buf: []rune(""), pos: 0},
			input:   []rune("aaa"),
			at:      0,
			want:    basicEditor{buf: []rune("aaa"), pos: 3},
		},
		{
			initial: basicEditor{buf: []rune("AAA"), pos: 2},
			input:   []rune("aaa"),
			at:      -1,
			want:    basicEditor{buf: []rune("aaaAAA"), pos: 5},
		},
		{
			initial: basicEditor{buf: []rune("AAA"), pos: 1},
			input:   []rune("aaa"),
			at:      2,
			want:    basicEditor{buf: []rune("AAaaaA"), pos: 1},
		},
		{
			initial: basicEditor{buf: []rune("AAA"), pos: 3},
			input:   []rune("aaa"),
			at:      10,
			want:    basicEditor{buf: []rune("AAAaaa"), pos: 6},
		},
	}
	for _, test := range tests {
		test.initial.insert(test.input, test.at)
		if !reflect.DeepEqual(test.initial, test.want) {
			t.Errorf("insert(%v, %v): got %v, want %v", test.input, test.at, test.initial, test.want)
		}
	}
}
