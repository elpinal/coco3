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
		{
			initial: basicEditor{buf: []rune("ABC"), pos: 1},
			input:   []rune("defg"),
			at:      1,
			want:    basicEditor{buf: []rune("AdefgBC"), pos: 5},
		},
	}
	for _, test := range tests {
		test.initial.insert(test.input, test.at)
		if !reflect.DeepEqual(test.initial, test.want) {
			t.Errorf("insert(%v, %v): got %v, want %v", test.input, test.at, test.initial, test.want)
		}
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		initial basicEditor
		from    int
		to      int
		want    basicEditor
	}{
		{
			initial: basicEditor{buf: []rune(""), pos: 0},
			from:    -1,
			to:      1,
			want:    basicEditor{buf: []rune(""), pos: 0},
		},
		{
			initial: basicEditor{buf: []rune("AAA"), pos: 2},
			from:    1,
			to:      2,
			want:    basicEditor{buf: []rune("AA"), pos: 1},
		},
		{
			initial: basicEditor{buf: []rune("AAA"), pos: 1},
			from:    2,
			to:      0,
			want:    basicEditor{buf: []rune("A"), pos: 0},
		},
		{
			initial: basicEditor{buf: []rune("AAAABBCCaaaa"), pos: 3},
			from:    4,
			to:      8,
			want:    basicEditor{buf: []rune("AAAAaaaa"), pos: 3},
		},
		{
			initial: basicEditor{buf: []rune(""), pos: 0},
			from:    -1,
			to:      -1,
			want:    basicEditor{buf: []rune(""), pos: 0},
		},
	}
	for _, test := range tests {
		test.initial.delete(test.from, test.to)
		if !reflect.DeepEqual(test.initial, test.want) {
			t.Errorf("delete(%v, %v): got %v, want %v", test.from, test.to, test.initial, test.want)
		}
	}
}

func TestSlice(t *testing.T) {
	tests := []struct {
		initial []rune
		from    int
		to      int
		want    []rune
	}{
		{
			initial: []rune(""),
			from:    -1,
			to:      1,
			want:    []rune(""),
		},
		{
			initial: []rune("ABC"),
			from:    -1,
			to:      1,
			want:    []rune("A"),
		},
		{
			initial: []rune("aaa bbb ccc"),
			from:    2,
			to:      9,
			want:    []rune("a bbb c"),
		},
		{
			initial: []rune("aaa x bbb"),
			from:    10,
			to:      4,
			want:    []rune("x bbb"),
		},
	}
	for _, test := range tests {
		e := &basicEditor{buf: test.initial}
		got := e.slice(test.from, test.to)
		if string(got) != string(test.want) {
			t.Errorf("slice(%v, %v): got %v, want %v", test.from, test.to, got, test.want)
		}
	}
}

func TestIndex(t *testing.T) {
	tests := []struct {
		initial []rune
		ch      rune
		start   int
		want    int
	}{
		{
			initial: []rune(""),
			ch:      'a',
			start:   0,
			want:    -1,
		},
		{
			initial: []rune("aaa"),
			ch:      'a',
			start:   1,
			want:    1,
		},
		{
			initial: []rune("abcde"),
			ch:      'e',
			start:   3,
			want:    4,
		},
		{
			initial: []rune("AA AA"),
			ch:      'A',
			start:   2,
			want:    3,
		},
	}
	for _, test := range tests {
		e := &basicEditor{buf: test.initial}
		got := e.index(test.ch, test.start)
		if got != test.want {
			t.Errorf("index(%v, %v): got %v, want %v", test.ch, test.start, got, test.want)
		}
	}
}

func TestLastIndex(t *testing.T) {
	tests := []struct {
		initial []rune
		ch      rune
		start   int
		want    int
	}{
		{
			initial: []rune(""),
			ch:      'a',
			start:   0,
			want:    -1,
		},
		{
			initial: []rune("aaa"),
			ch:      'a',
			start:   1,
			want:    0,
		},
		{
			initial: []rune("abcde"),
			ch:      'e',
			start:   3,
			want:    -1,
		},
		{
			initial: []rune("AA AA"),
			ch:      'A',
			start:   2,
			want:    1,
		},
	}
	for _, test := range tests {
		e := &basicEditor{buf: test.initial}
		got := e.lastIndex(test.ch, test.start)
		if got != test.want {
			t.Errorf("lastIndex(%v, %v): got %v, want %v", string(test.ch), test.start, got, test.want)
		}
	}
}

func TestReplace(t *testing.T) {
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
			want:    basicEditor{buf: []rune("aaa"), pos: 0},
		},
		{
			initial: basicEditor{buf: []rune("AAA"), pos: 2},
			input:   []rune("aaa"),
			at:      -1,
			want:    basicEditor{buf: []rune("aaA"), pos: 2},
		},
		{
			initial: basicEditor{buf: []rune("AAA"), pos: 1},
			input:   []rune("aaa"),
			at:      2,
			want:    basicEditor{buf: []rune("AAaaa"), pos: 1},
		},
		{
			initial: basicEditor{buf: []rune("AAA"), pos: 3},
			input:   []rune("aaa"),
			at:      10,
			want:    basicEditor{buf: []rune("AAA       aaa"), pos: 3},
		},
	}
	for _, test := range tests {
		test.initial.replace(test.input, test.at)
		if !reflect.DeepEqual(test.initial, test.want) {
			t.Errorf("replace(%v, %v): got %v, want %v", string(test.input), test.at, test.initial, test.want)
		}
	}
}