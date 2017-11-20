package commandline

import (
	"reflect"
	"testing"
)

func TestScan(t *testing.T) {
	tests := []struct {
		src  string
		want []token
	}{
		{
			src: "",
			want: []token{
				{
					tt:    tokenEOF,
					value: []byte(""),
				},
			},
		},
		{
			src: " ",
			want: []token{
				{
					tt:    tokenEOF,
					value: []byte(""),
				},
			},
		},
		{
			src: "a",
			want: []token{
				{
					tt:    ident,
					value: []byte("a"),
				},
			},
		},
		{
			src: "abc",
			want: []token{
				{
					tt:    ident,
					value: []byte("abc"),
				},
			},
		},
		{
			src: "ABc",
			want: []token{
				{
					tt:    ident,
					value: []byte("ABc"),
				},
			},
		},
		{
			src: " abc",
			want: []token{
				{
					tt:    ident,
					value: []byte("abc"),
				},
			},
		},
		{
			src: "a b",
			want: []token{
				{
					tt:    ident,
					value: []byte("a"),
				},
				{
					tt:    ident,
					value: []byte("b"),
				},
			},
		},
		{
			src: " a b",
			want: []token{
				{
					tt:    ident,
					value: []byte("a"),
				},
				{
					tt:    ident,
					value: []byte("b"),
				},
			},
		},
		{
			src: "a    b cd  e fgh",
			want: []token{
				{
					tt:    ident,
					value: []byte("a"),
				},
				{
					tt:    ident,
					value: []byte("b"),
				},
				{
					tt:    ident,
					value: []byte("cd"),
				},
				{
					tt:    ident,
					value: []byte("e"),
				},
				{
					tt:    ident,
					value: []byte("fgh"),
				},
			},
		},
		{
			src: "1",
			want: []token{
				{
					tt:    tokenErr,
					value: []byte("1"),
				},
			},
		},
		{
			src: `"a"`,
			want: []token{
				{
					tt:    str,
					value: []byte(`"a"`),
				},
			},
		},
		{
			src: `a "b"`,
			want: []token{
				{
					tt:    ident,
					value: []byte("a"),
				},
				{
					tt:    str,
					value: []byte(`"b"`),
				},
			},
		},
		{
			src: `a"b"`,
			want: []token{
				{
					tt:    ident,
					value: []byte("a"),
				},
				{
					tt:    str,
					value: []byte(`"b"`),
				},
			},
		},
		{
			src: `a "b`,
			want: []token{
				{
					tt:    ident,
					value: []byte("a"),
				},
				{
					tt:    tokenErr,
					value: []byte(`"b`),
				},
			},
		},
	}
	for i, test := range tests {
		s := scan([]byte(test.src))
		for _, want := range test.want {
			got := <-s.tokens
			if !reflect.DeepEqual(got, want) {
				t.Errorf("scan/%d: got %v, want %v", i, got, want)
			}
		}
	}
}
