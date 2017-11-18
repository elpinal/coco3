package editor

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
		e := &basic{buf: s}
		e.move(test.input)
		if e.pos != test.want {
			t.Errorf("move(%v): got %v, want %v", test.input, e.pos, test.want)
		}
	}
}

func TestInsert(t *testing.T) {
	tests := []struct {
		initial basic
		input   []rune
		at      int
		want    basic
	}{
		{
			initial: basic{buf: []rune(""), pos: 0},
			input:   []rune("aaa"),
			at:      0,
			want:    basic{buf: []rune("aaa"), pos: 3},
		},
		{
			initial: basic{buf: []rune("AAA"), pos: 2},
			input:   []rune("aaa"),
			at:      -1,
			want:    basic{buf: []rune("aaaAAA"), pos: 5},
		},
		{
			initial: basic{buf: []rune("AAA"), pos: 1},
			input:   []rune("aaa"),
			at:      2,
			want:    basic{buf: []rune("AAaaaA"), pos: 1},
		},
		{
			initial: basic{buf: []rune("AAA"), pos: 3},
			input:   []rune("aaa"),
			at:      10,
			want:    basic{buf: []rune("AAAaaa"), pos: 6},
		},
		{
			initial: basic{buf: []rune("ABC"), pos: 1},
			input:   []rune("defg"),
			at:      1,
			want:    basic{buf: []rune("AdefgBC"), pos: 5},
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
		initial basic
		from    int
		to      int
		want    basic
	}{
		{
			initial: basic{buf: []rune(""), pos: 0},
			from:    -1,
			to:      1,
			want:    basic{buf: []rune(""), pos: 0},
		},
		{
			initial: basic{buf: []rune("AAA"), pos: 2},
			from:    1,
			to:      2,
			want:    basic{buf: []rune("AA"), pos: 1},
		},
		{
			initial: basic{buf: []rune("AAA"), pos: 1},
			from:    2,
			to:      0,
			want:    basic{buf: []rune("A"), pos: 0},
		},
		{
			initial: basic{buf: []rune("AAAABBCCaaaa"), pos: 3},
			from:    4,
			to:      8,
			want:    basic{buf: []rune("AAAAaaaa"), pos: 3},
		},
		{
			initial: basic{buf: []rune(""), pos: 0},
			from:    -1,
			to:      -1,
			want:    basic{buf: []rune(""), pos: 0},
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
		e := &basic{buf: test.initial}
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
		e := &basic{buf: test.initial}
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
		e := &basic{buf: test.initial}
		got := e.lastIndex(test.ch, test.start)
		if got != test.want {
			t.Errorf("lastIndex(%v, %v): got %v, want %v", string(test.ch), test.start, got, test.want)
		}
	}
}

func TestReplace(t *testing.T) {
	tests := []struct {
		initial basic
		input   []rune
		at      int
		want    basic
	}{
		{
			initial: basic{buf: []rune(""), pos: 0},
			input:   []rune("aaa"),
			at:      0,
			want:    basic{buf: []rune("aaa"), pos: 0},
		},
		{
			initial: basic{buf: []rune("AAA"), pos: 2},
			input:   []rune("aaa"),
			at:      -1,
			want:    basic{buf: []rune("aaA"), pos: 2},
		},
		{
			initial: basic{buf: []rune("AAA"), pos: 1},
			input:   []rune("aaa"),
			at:      2,
			want:    basic{buf: []rune("AAaaa"), pos: 1},
		},
		{
			initial: basic{buf: []rune("AAA"), pos: 3},
			input:   []rune("aaa"),
			at:      10,
			want:    basic{buf: []rune("AAA       aaa"), pos: 3},
		},
	}
	for _, test := range tests {
		test.initial.replace(test.input, test.at)
		if !reflect.DeepEqual(test.initial, test.want) {
			t.Errorf("replace(%v, %v): got %v, want %v", string(test.input), test.at, test.initial, test.want)
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
			initial: basic{buf: []rune("aaa "), pos: 3},
			want:    4,
		},
		{
			initial: basic{buf: []rune("().aaa"), pos: 1},
			want:    3,
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
			initial: basic{buf: []rune(" aaa"), pos: 1},
			want:    0,
		},
		{
			initial: basic{buf: []rune("   "), pos: 2},
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
		{
			initial: basic{buf: []rune("     "), pos: 0},
			want:    5,
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

func TestWordEnd(t *testing.T) {
	tests := []struct {
		initial basic
		want    int
	}{
		{
			initial: basic{buf: []rune(""), pos: 0},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa "), pos: 3},
			want:    4,
		},
		{
			initial: basic{buf: []rune("().aaa"), pos: 1},
			want:    2,
		},
		{
			initial: basic{buf: []rune("aaa"), pos: 0},
			want:    2,
		},
		{
			initial: basic{buf: []rune("aaa()"), pos: 2},
			want:    4,
		},
		{
			initial: basic{buf: []rune("aaa x bbb"), pos: 3},
			want:    4,
		},
		{
			initial: basic{buf: []rune("aaa () bbb"), pos: 3},
			want:    5,
		},
		{
			initial: basic{buf: []rune("##### x bbb"), pos: 3},
			want:    4,
		},
		{
			initial: basic{buf: []rune("#####   aa#"), pos: 5},
			want:    9,
		},
		{
			initial: basic{buf: []rune("  ###aaa"), pos: 0},
			want:    4,
		},
	}
	for i, test := range tests {
		e := &editor{basic: test.initial}
		e.wordEnd()
		if e.pos != test.want {
			t.Errorf("wordEnd %v: got %v, want %v", i, e.pos, test.want)
		}
	}
}

func TestWordEndNonBlank(t *testing.T) {
	tests := []struct {
		initial basic
		want    int
	}{
		{
			initial: basic{buf: []rune(""), pos: 0},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa "), pos: 3},
			want:    4,
		},
		{
			initial: basic{buf: []rune("().aaa"), pos: 1},
			want:    5,
		},
		{
			initial: basic{buf: []rune("aaa"), pos: 0},
			want:    2,
		},
		{
			initial: basic{buf: []rune("aaa()"), pos: 2},
			want:    4,
		},
		{
			initial: basic{buf: []rune("aaa x bbb"), pos: 3},
			want:    4,
		},
		{
			initial: basic{buf: []rune("aaa () bbb"), pos: 3},
			want:    5,
		},
		{
			initial: basic{buf: []rune("##### x bbb"), pos: 3},
			want:    4,
		},
		{
			initial: basic{buf: []rune("#####   aa#"), pos: 5},
			want:    10,
		},
		{
			initial: basic{buf: []rune("  # "), pos: 0},
			want:    2,
		},
	}
	for i, test := range tests {
		e := &editor{basic: test.initial}
		e.wordEndNonBlank()
		if e.pos != test.want {
			t.Errorf("wordEndNonBlank %v: got %v, want %v", i, e.pos, test.want)
		}
	}
}

func TestWordEndBackward(t *testing.T) {
	tests := []struct {
		initial basic
		want    int
	}{
		{
			initial: basic{buf: []rune(""), pos: 0},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa "), pos: 3},
			want:    2,
		},
		{
			initial: basic{buf: []rune("aaa.()"), pos: 4},
			want:    2,
		},
		{
			initial: basic{buf: []rune("aaa"), pos: 2},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa()"), pos: 3},
			want:    2,
		},
		{
			initial: basic{buf: []rune("aaa x bbb"), pos: 5},
			want:    4,
		},
		{
			initial: basic{buf: []rune("aaa () bbb"), pos: 6},
			want:    5,
		},
		{
			initial: basic{buf: []rune("##### x bbb"), pos: 9},
			want:    6,
		},
		{
			initial: basic{buf: []rune("#####   aa#"), pos: 9},
			want:    4,
		},
		{
			initial: basic{buf: []rune("  # "), pos: 3},
			want:    2,
		},
		{
			initial: basic{buf: []rune("   "), pos: 1},
			want:    0,
		},
		{
			initial: basic{buf: []rune(" %a "), pos: 2},
			want:    1,
		},
		{
			initial: basic{buf: []rune(" a % "), pos: 3},
			want:    1,
		},
	}
	for i, test := range tests {
		e := &editor{basic: test.initial}
		e.wordEndBackward()
		if e.pos != test.want {
			t.Errorf("wordEndBackward %v: got %v, want %v", i, e.pos, test.want)
		}
	}
}

func TestWordEndBackwardNonBlank(t *testing.T) {
	tests := []struct {
		initial basic
		want    int
	}{
		{
			initial: basic{buf: []rune(""), pos: 0},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa "), pos: 3},
			want:    2,
		},
		{
			initial: basic{buf: []rune("aaa.()"), pos: 4},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa"), pos: 2},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa()"), pos: 1},
			want:    0,
		},
		{
			initial: basic{buf: []rune("aaa x bbb"), pos: 5},
			want:    4,
		},
		{
			initial: basic{buf: []rune("aaa () bbb"), pos: 6},
			want:    5,
		},
		{
			initial: basic{buf: []rune("##### x bbb"), pos: 9},
			want:    6,
		},
		{
			initial: basic{buf: []rune("#####   aa#"), pos: 9},
			want:    4,
		},
		{
			initial: basic{buf: []rune("  # "), pos: 3},
			want:    2,
		},
	}
	for i, test := range tests {
		e := &editor{basic: test.initial}
		e.wordEndBackwardNonBlank()
		if e.pos != test.want {
			t.Errorf("wordEndBackwardNonBlank %v: got %v, want %v", i, e.pos, test.want)
		}
	}
}
