package editor

import (
	"testing"

	"github.com/elpinal/coco3/editor/register"
)

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
			r:    register.Unnamed,
			at:   0,
			want: []rune(""),
		},
		{
			buf:  []rune("ABCDE"),
			from: 2,
			to:   5,
			r:    register.Unnamed,
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
			r:    register.BlackHole,
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
		r := register.Registers{}
		r.Init()
		e := &editor{basic: basic{buf: test.buf}, Registers: r}
		e.yank(test.r, test.from, test.to)
		e.put(test.r, test.at)
		if string(e.buf) != string(test.want) {
			t.Errorf("%v: got %v, want %v", i, string(e.buf), string(test.want))
		}
	}
}

func TestWordForward(t *testing.T) {
	tests := []struct {
		initial basic
		want    int
	}{
		{
			initial: basic{buf: []rune(""), pos: 0},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa"), pos: 0},
			want:    3,
		},
		{
			initial: basic{buf: []rune("aaa()"), pos: 2},
			want:    3,
		},
		{
			initial: basic{buf: []rune("aaa x bbb"), pos: 3},
			want:    4,
		},
		{
			initial: basic{buf: []rune("aaa () bbb"), pos: 3},
			want:    4,
		},
		{
			initial: basic{buf: []rune("##### x bbb"), pos: 3},
			want:    6,
		},
		{
			initial: basic{buf: []rune("#####   aa#"), pos: 5},
			want:    8,
		},
	}
	for i, test := range tests {
		e := &editor{basic: test.initial}
		e.wordForward()
		if e.pos != test.want {
			t.Errorf("wordForward %v: got %v, want %v", i, e.pos, test.want)
		}
	}
}

func TestWordBackward(t *testing.T) {
	tests := []struct {
		initial basic
		want    int
	}{
		{
			initial: basic{buf: []rune(""), pos: 0},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa"), pos: 3},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa()"), pos: 4},
			want:    3,
		},
		{
			initial: basic{buf: []rune("aaa x bbb"), pos: 5},
			want:    4,
		},
		{
			initial: basic{buf: []rune("aaa () bbb"), pos: 5},
			want:    4,
		},
		{
			initial: basic{buf: []rune("aaa x #####"), pos: 8},
			want:    6,
		},
		{
			initial: basic{buf: []rune("#aa   #####"), pos: 5},
			want:    1,
		},
	}
	for i, test := range tests {
		e := &editor{basic: test.initial}
		e.wordBackward()
		if e.pos != test.want {
			t.Errorf("wordBackward %v: got %v, want %v", i, e.pos, test.want)
		}
	}
}

func TestWordForwardNonBlank(t *testing.T) {
	tests := []struct {
		initial basic
		want    int
	}{
		{
			initial: basic{buf: []rune(""), pos: 0},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa"), pos: 0},
			want:    3,
		},
		{
			initial: basic{buf: []rune("aaa()"), pos: 2},
			want:    5,
		},
		{
			initial: basic{buf: []rune("aaa x bbb"), pos: 3},
			want:    4,
		},
		{
			initial: basic{buf: []rune("aaa () bbb"), pos: 3},
			want:    4,
		},
		{
			initial: basic{buf: []rune("##### x bbb"), pos: 3},
			want:    6,
		},
		{
			initial: basic{buf: []rune("#####   aa#"), pos: 5},
			want:    8,
		},
	}
	for i, test := range tests {
		e := &editor{basic: test.initial}
		e.wordForwardNonBlank()
		if e.pos != test.want {
			t.Errorf("wordForwardNonBlank %v: got %v, want %v", i, e.pos, test.want)
		}
	}
}

func TestWordBackwardNonBlank(t *testing.T) {
	tests := []struct {
		initial basic
		want    int
	}{
		{
			initial: basic{buf: []rune(""), pos: 0},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa"), pos: 3},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa()"), pos: 4},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa x bbb"), pos: 5},
			want:    4,
		},
		{
			initial: basic{buf: []rune("aaa () bbb"), pos: 5},
			want:    4,
		},
		{
			initial: basic{buf: []rune("aaa x #####"), pos: 8},
			want:    6,
		},
		{
			initial: basic{buf: []rune("#aa   #####"), pos: 5},
			want:    0,
		},
	}
	for i, test := range tests {
		e := &editor{basic: test.initial}
		e.wordBackwardNonBlank()
		if e.pos != test.want {
			t.Errorf("wordBackwardNonBlank %v: got %v, want %v", i, e.pos, test.want)
		}
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		input []rune
		from  int
		to    int
		want  []rune
	}{
		{
			input: []rune(""),
			from:  0,
			to:    0,
			want:  []rune(""),
		},
		{
			input: []rune("Gopher"),
			from:  0,
			to:    8,
			want:  []rune("GOPHER"),
		},
		{
			input: []rune("AAAAAA"),
			from:  -9,
			to:    9,
			want:  []rune("AAAAAA"),
		},
		{
			input: []rune("aaa X bbb X ccc"),
			from:  4,
			to:    8,
			want:  []rune("aaa X BBb X ccc"),
		},
	}
	for i, test := range tests {
		e := &editor{basic: basic{buf: test.input}}
		e.toUpper(test.from, test.to)
		if string(e.buf) != string(test.want) {
			t.Errorf("toUpper %v: got %v, want %v", i, string(e.buf), string(test.want))
		}
	}
}

func TestToLower(t *testing.T) {
	tests := []struct {
		input []rune
		from  int
		to    int
		want  []rune
	}{
		{
			input: []rune(""),
			from:  0,
			to:    0,
			want:  []rune(""),
		},
		{
			input: []rune("Gopher"),
			from:  0,
			to:    8,
			want:  []rune("gopher"),
		},
		{
			input: []rune("AAAAAA"),
			from:  -9,
			to:    9,
			want:  []rune("aaaaaa"),
		},
		{
			input: []rune("aaa X bbb X ccc"),
			from:  4,
			to:    8,
			want:  []rune("aaa x bbb X ccc"),
		},
	}
	for i, test := range tests {
		e := &editor{basic: basic{buf: test.input}}
		e.toLower(test.from, test.to)
		if string(e.buf) != string(test.want) {
			t.Errorf("toLower %v: got %v, want %v", i, string(e.buf), string(test.want))
		}
	}
}

func TestCurrentWord(t *testing.T) {
	tests := []struct {
		input   []rune
		pos     int
		include bool
		from    int
		to      int
	}{
		{
			input:   []rune(""),
			pos:     0,
			include: false,
			from:    0,
			to:      0,
		},
		{
			input:   []rune("aaa"),
			pos:     0,
			include: false,
			from:    0,
			to:      3,
		},
		{
			input:   []rune("a a a"),
			pos:     1,
			include: false,
			from:    1,
			to:      2,
		},
		{
			input:   []rune(" aaa bbb ccc "),
			pos:     7,
			include: true,
			from:    5,
			to:      9,
		},
	}
	for i, test := range tests {
		e := &editor{basic: basic{buf: test.input, pos: test.pos}}
		from, to := e.currentWord(test.include)
		if from != test.from {
			t.Errorf("currentWord/%v (from): got %v, want %v", i, from, test.from)
		}
		if to != test.to {
			t.Errorf("currentWord/%v (to): got %v, want %v", i, to, test.to)
		}
	}
}
