package main

import "testing"

func TestRegister(t *testing.T) {
	tests := []struct {
		buf  []rune
		from int
		to   int
		r    rune
		at   int
		want []rune
	}{
		{
			buf:  []rune(""),
			from: -1,
			to:   1,
			r:    registerUnnamed,
			at:   0,
			want: []rune(""),
		},
		{
			buf:  []rune("ABCDE"),
			from: 2,
			to:   5,
			r:    registerUnnamed,
			at:   1,
			want: []rune("ACDEBCDE"),
		},
		{
			buf:  []rune("A B C"),
			from: 0,
			to:   5,
			r:    '5',
			at:   4,
			want: []rune("A B A B CC"),
		},
		{
			buf:  []rune("A B C"),
			from: 0,
			to:   5,
			r:    registerBlackHole,
			at:   4,
			want: []rune("A B C"),
		},
		{
			buf:  []rune("A"),
			from: 0,
			to:   1,
			r:    'A',
			at:   1,
			want: []rune("AA"),
		},
	}
	for i, test := range tests {
		r := registers{}
		r.init()
		e := &editor{basicEditor: basicEditor{buf: test.buf}, registers: r}
		e.yank(test.r, test.from, test.to)
		e.put(test.r, test.at)
		if string(e.buf) != string(test.want) {
			t.Errorf("%v: got %v, want %v", i, string(e.buf), string(test.want))
		}
	}
}

func TestWordForward(t *testing.T) {
	tests := []struct {
		initial basicEditor
		want int
	}{
		{
			initial: basicEditor{buf: []rune(""), pos: 0},
			want: 0,
		},
		{
			initial: basicEditor{buf: []rune("aaa"), pos: 0},
			want: 3,
		},
		{
			initial: basicEditor{buf: []rune("aaa()"), pos: 2},
			want: 3,
		},
		{
			initial: basicEditor{buf: []rune("aaa x bbb"), pos: 3},
			want: 4,
		},
		{
			initial: basicEditor{buf: []rune("aaa () bbb"), pos: 3},
			want: 4,
		},
		{
			initial: basicEditor{buf: []rune("##### x bbb"), pos: 3},
			want: 6,
		},
	}
	for i, test := range tests {
		e := &editor{basicEditor: test.initial}
		e.wordForward()
		if e.pos != test.want {
			t.Errorf("wordForward %v: got %v, want %v", i, e.pos, test.want)
		}
	}
}
