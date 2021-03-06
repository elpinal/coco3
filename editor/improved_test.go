package editor

import (
	"reflect"
	"strings"
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

func BenchmarkToUpper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := editor{basic: basic{buf: []rune("aaa BBB ccc")}}
		e.toUpper(0, 11)
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

func BenchmarkToLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := editor{basic: basic{buf: []rune("aaa BBB ccc")}}
		e.toLower(0, 11)
	}
}

func TestSwitchCase(t *testing.T) {
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
			want:  []rune("gOPHER"),
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
			want:  []rune("aaa x BBb X ccc"),
		},
	}
	for i, test := range tests {
		e := &editor{basic: basic{buf: test.input}}
		e.switchCase(test.from, test.to)
		if string(e.buf) != string(test.want) {
			t.Errorf("switchCase %v: got %v, want %v", i, string(e.buf), string(test.want))
		}
	}
}

func BenchmarkSwitchCase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := editor{basic: basic{buf: []rune("aaa BBB ccc")}}
		e.switchCase(0, 11)
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
		{
			input:   []rune("a"),
			pos:     1,
			include: true,
			from:    -1,
			to:      -1,
		},
		{
			input:   []rune("%"),
			pos:     0,
			include: true,
			from:    0,
			to:      1,
		},
		{
			input:   []rune(" %"),
			pos:     1,
			include: true,
			from:    0,
			to:      2,
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

func TestCurrentWordNonBlank(t *testing.T) {
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
			input:   []rune("a%a"),
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
			input:   []rune(" aaa .bb ccc "),
			pos:     7,
			include: true,
			from:    5,
			to:      9,
		},
		{
			input:   []rune("aa"),
			pos:     2,
			include: true,
			from:    -1,
			to:      -1,
		},
		{
			input:   []rune(" aa"),
			pos:     2,
			include: true,
			from:    0,
			to:      3,
		},
	}
	for i, test := range tests {
		e := &editor{basic: basic{buf: test.input, pos: test.pos}}
		from, to := e.currentWordNonBlank(test.include)
		if from != test.from {
			t.Errorf("currentWordNonBlank/%v (from): got %v, want %v", i, from, test.from)
		}
		if to != test.to {
			t.Errorf("currentWordNonBlank/%v (to): got %v, want %v", i, to, test.to)
		}
	}
}

func TestCurrentQuote(t *testing.T) {
	tests := []struct {
		input   []rune
		pos     int
		include bool
		quote   rune
		from    int
		to      int
	}{
		{
			input:   []rune(""),
			pos:     0,
			include: false,
			quote:   '\'',
			from:    0,
			to:      0,
		},
		{
			input:   []rune("'aaa'"),
			pos:     0,
			include: false,
			quote:   '\'',
			from:    1,
			to:      4,
		},
		{
			input:   []rune("a' a 'a"),
			pos:     3,
			include: false,
			quote:   '\'',
			from:    2,
			to:      5,
		},
		{
			input:   []rune(` aaa "bbb ccc "`),
			pos:     7,
			include: true,
			quote:   '"',
			from:    4,
			to:      15,
		},
		{
			input:   []rune(` aaa "bbb ccc "ddd "eee f"ff`),
			pos:     25,
			include: true,
			quote:   '"',
			from:    18,
			to:      26,
		},
		{
			input:   []rune(" aaa `bbb ccc `"),
			pos:     7,
			include: true,
			quote:   '`',
			from:    4,
			to:      15,
		},
		{
			input:   []rune("abc"),
			pos:     3,
			include: true,
			quote:   '`',
			from:    -1,
			to:      -1,
		},
		{
			input:   []rune(" abc` "),
			pos:     2,
			include: true,
			quote:   '`',
			from:    -1,
			to:      -1,
		},
		{
			input:   []rune(" `abc "),
			pos:     3,
			include: true,
			quote:   '`',
			from:    -1,
			to:      -1,
		},
		{
			input:   []rune(" `abc` "),
			pos:     3,
			include: true,
			quote:   '`',
			from:    1,
			to:      7,
		},
	}
	for i, test := range tests {
		e := &editor{basic: basic{buf: test.input, pos: test.pos}}
		from, to := e.currentQuote(test.include, test.quote)
		if from != test.from {
			t.Errorf("currentQuote/%v (from): got %v, want %v", i, from, test.from)
		}
		if to != test.to {
			t.Errorf("currentQuote/%v (to): got %v, want %v", i, to, test.to)
		}
	}
}

func TestCurrentParen(t *testing.T) {
	tests := []struct {
		input   []rune
		pos     int
		include bool
		lparen  rune
		rparen  rune
		from    int
		to      int
	}{
		{
			input:   []rune(""),
			pos:     0,
			include: false,
			lparen:  '(',
			rparen:  ')',
			from:    0,
			to:      0,
		},
		{
			input:   []rune("(aaa)"),
			pos:     0,
			include: false,
			lparen:  '(',
			rparen:  ')',
			from:    1,
			to:      4,
		},
		{
			input:   []rune("a( a )a"),
			pos:     3,
			include: false,
			lparen:  '(',
			rparen:  ')',
			from:    2,
			to:      5,
		},
		{
			input:   []rune(` aaa <bbb ccc >`),
			pos:     7,
			include: true,
			lparen:  '<',
			rparen:  '>',
			from:    5,
			to:      15,
		},
		{
			input:   []rune(` aaa "bbb ccc "ddd {eee f}ff`),
			pos:     25,
			include: true,
			lparen:  '{',
			rparen:  '}',
			from:    19,
			to:      26,
		},
		{
			input:   []rune(" aaa [bbb ccc ]"),
			pos:     7,
			include: true,
			lparen:  '[',
			rparen:  ']',
			from:    5,
			to:      15,
		},
		{
			input:   []rune("a"),
			pos:     1,
			include: true,
			lparen:  '[',
			rparen:  ']',
			from:    -1,
			to:      -1,
		},
		{
			input:   []rune("a]"),
			pos:     1,
			include: true,
			lparen:  '[',
			rparen:  ']',
			from:    -1,
			to:      -1,
		},
		{
			input:   []rune("[a"),
			pos:     1,
			include: true,
			lparen:  '[',
			rparen:  ']',
			from:    -1,
			to:      -1,
		},
	}
	for i, test := range tests {
		e := &editor{basic: basic{buf: test.input, pos: test.pos}}
		from, to := e.currentParen(test.include, test.lparen, test.rparen)
		if from != test.from {
			t.Errorf("currentQuote/%v (from): got %v, want %v", i, from, test.from)
		}
		if to != test.to {
			t.Errorf("currentQuote/%v (to): got %v, want %v", i, to, test.to)
		}
	}
}

func TestSearchLeft(t *testing.T) {
	tests := []struct {
		input  []rune
		pos    int
		lparen rune
		rparen rune
		want   int
	}{
		{
			input:  []rune(""),
			pos:    0,
			lparen: '(',
			rparen: ')',
			want:   -1,
		},
		{
			input:  []rune("()"),
			pos:    0,
			lparen: '(',
			rparen: ')',
			want:   0,
		},
		{
			input:  []rune("a(a(a)a)a"),
			pos:    6,
			lparen: '(',
			rparen: ')',
			want:   1,
		},
		{
			input:  []rune("( () )"),
			pos:    2,
			lparen: '(',
			rparen: ')',
			want:   2,
		},
		{
			input:  []rune("abababab"),
			pos:    3,
			lparen: 'a',
			rparen: 'b',
			want:   2,
		},
	}
	for i, test := range tests {
		e := &editor{basic: basic{buf: test.input, pos: test.pos}}
		got := e.searchLeft(test.lparen, test.rparen)
		if got != test.want {
			t.Errorf("searchLeft/%v: got %v, want %v", i, got, test.want)
		}
	}
}

func TestSearchRight(t *testing.T) {
	tests := []struct {
		input  []rune
		pos    int
		lparen rune
		rparen rune
		want   int
	}{
		{
			input:  []rune(""),
			pos:    0,
			lparen: '(',
			rparen: ')',
			want:   -1,
		},
		{
			input:  []rune("()"),
			pos:    0,
			lparen: '(',
			rparen: ')',
			want:   1,
		},
		{
			input:  []rune("a(a(a)a)a"),
			pos:    2,
			lparen: '(',
			rparen: ')',
			want:   7,
		},
		{
			input:  []rune("( () )"),
			pos:    3,
			lparen: '(',
			rparen: ')',
			want:   3,
		},
		{
			input:  []rune("abababab"),
			pos:    4,
			lparen: 'a',
			rparen: 'b',
			want:   5,
		},
	}
	for i, test := range tests {
		e := &editor{basic: basic{buf: test.input, pos: test.pos}}
		got := e.searchRight(test.lparen, test.rparen)
		if got != test.want {
			t.Errorf("searchRight/%v: got %v, want %v", i, got, test.want)
		}
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		buf   string
		pos   int
		query string
		found bool
		sr    searchRange
	}{
		{
			buf:   "",
			pos:   0,
			query: "",
			found: false,
			sr:    nil,
		},
		{
			buf:   "",
			pos:   0,
			query: "aaa",
			found: false,
			sr:    nil,
		},
		{
			buf:   "aaabbb",
			pos:   0,
			query: "a",
			found: true,
			sr: [][]int{
				{0, 1},
				{1, 2},
				{2, 3},
			},
		},
		{
			buf:   "aaabbb xy cccddd xy abcd",
			pos:   10,
			query: "xy",
			found: true,
			sr: [][]int{
				{7, 9},
				{17, 19},
			},
		},
	}
	for i, test := range tests {
		e := &editor{basic: basic{buf: []rune(test.buf), pos: test.pos}}
		found := e.search(test.query)
		if found != test.found {
			t.Errorf("search/%v (found): got %v, want %v", i, found, test.found)
		}
		if !reflect.DeepEqual(e.sr, test.sr) {
			t.Errorf("search/%v (searchRange): got %v, want %v", i, e.sr, test.sr)
		}
	}
}

func BenchmarkSearchLeftSimple(b *testing.B) {
	e := newEditor()
	e.insert([]rune("("+strings.Repeat(" ", 199)), 0)
	e.move(199)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.searchLeft('(', ')')
	}
}

func BenchmarkSearchLeft(b *testing.B) {
	e := newEditor()
	e.insert([]rune(strings.Repeat("(", 100)+strings.Repeat(")", 100)), 0)
	e.move(199)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.searchLeft('(', ')')
	}
}

func BenchmarkSearchRightSimple(b *testing.B) {
	e := newEditor()
	e.insert([]rune(strings.Repeat(" ", 199)+")"), 0)
	e.move(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.searchRight('(', ')')
	}
}

func BenchmarkSearchRight(b *testing.B) {
	e := newEditor()
	e.insert([]rune(strings.Repeat("(", 100)+strings.Repeat(")", 100)), 0)
	e.move(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.searchRight('(', ')')
	}
}

func TestCharSearch(t *testing.T) {
	tests := []struct {
		input []rune
		pos   int
		char  rune
		want  int
		fail  bool
	}{
		{
			input: []rune("aaa bcd eee fgh"),
			pos:   0,
			char:  'd',
			want:  6,
		},
		{
			input: []rune("aaa bcd eee fgh"),
			pos:   0,
			char:  'X',
			fail:  true,
		},
	}
	for i, test := range tests {
		e := newEditorBuffer([]rune(test.input))
		e.move(test.pos)
		n, err := e.charSearch(test.char)
		if !test.fail && err != nil {
			t.Errorf("charSearch/%d: %v", i, err)
		}
		if n != test.want {
			t.Errorf("charSearch/%d: want %d, got %d", i, test.want, n)
		}
	}
}

func BenchmarkCharSearch(b *testing.B) {
	e := newEditorBuffer([]rune("aaa bcd eee fgh"))
	for i := 0; i < b.N; i++ {
		_, _ = e.charSearch('d')
	}
}

func TestCharSearchBefore(t *testing.T) {
	e := newEditorBuffer([]rune("aaa bcd eee fgh"))
	i, err := e.charSearchBefore('d')
	if err != nil {
		t.Errorf("charSearchBefore: %v", err)
	}
	if want := 5; i != want {
		t.Errorf("charSearchBefore: want %d, got %d", want, i)
	}
}

func BenchmarkCharSearchBefore(b *testing.B) {
	e := newEditorBuffer([]rune("aaa bcd eee fgh"))
	for i := 0; i < b.N; i++ {
		_, _ = e.charSearchBefore('d')
	}
}

func TestCharSearchBackward(t *testing.T) {
	e := newEditorBuffer([]rune("aaa bcd eee fgh"))
	e.move(10)
	i, err := e.charSearchBackward('b')
	if err != nil {
		t.Errorf("charSearchBackward: %v", err)
	}
	if want := 4; i != want {
		t.Errorf("charSearchBackward: want %d, got %d", want, i)
	}
}

func BenchmarkCharSearchBackward(b *testing.B) {
	e := newEditorBuffer([]rune("aaa bcd eee fgh"))
	e.move(10)
	for i := 0; i < b.N; i++ {
		_, _ = e.charSearchBackward('b')
	}
}

func TestCharSearchBackwardAfter(t *testing.T) {
	e := newEditorBuffer([]rune("aaa bcd eee fgh"))
	e.move(10)
	i, err := e.charSearchBackwardAfter('b')
	if err != nil {
		t.Errorf("charSearchBackwardAfter: %v", err)
	}
	if want := 5; i != want {
		t.Errorf("charSearchBackwardAfter: want %d, got %d", want, i)
	}
}

func BenchmarkCharSearchBackwardAfter(b *testing.B) {
	e := newEditorBuffer([]rune("aaa bcd eee fgh"))
	e.move(10)
	for i := 0; i < b.N; i++ {
		_, _ = e.charSearchBackwardAfter('b')
	}
}

func TestOverwrite(t *testing.T) {
	tests := []struct {
		base  []rune
		cover []rune
		at    int
		want  string
	}{
		{
			base:  []rune("abcabc"),
			cover: []rune("123"),
			at:    0,
			want:  "123abc",
		},
		{
			base:  []rune(""),
			cover: []rune("123"),
			at:    0,
			want:  "123",
		},
	}
	for i, test := range tests {
		e := newEditor()
		s := string(e.overwrite(test.base, test.cover, test.at))
		if s != test.want {
			t.Errorf("charSearch/%d: want %q, got %q", i, test.want, s)
		}
	}
}

func BenchmarkOverwrite(b *testing.B) {
	e := newEditor()
	for i := 0; i < b.N; i++ {
		_ = e.overwrite([]rune("1234567890"), []rune("0987654321"), 1)
	}
}
