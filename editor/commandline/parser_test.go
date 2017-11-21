package commandline

import (
	"reflect"
	"testing"
)

func TestScan(t *testing.T) {
	tests := []struct {
		src  string
		want []Token
	}{
		{
			src: "",
			want: []Token{
				{
					Type:  tokenEOF,
					Value: []byte(""),
				},
			},
		},
		{
			src: " ",
			want: []Token{
				{
					Type:  tokenEOF,
					Value: []byte(""),
				},
			},
		},
		{
			src: "a",
			want: []Token{
				{
					Type:  tokenIdent,
					Value: []byte("a"),
				},
			},
		},
		{
			src: "abc",
			want: []Token{
				{
					Type:  tokenIdent,
					Value: []byte("abc"),
				},
			},
		},
		{
			src: "ABc",
			want: []Token{
				{
					Type:  tokenIdent,
					Value: []byte("ABc"),
				},
			},
		},
		{
			src: " abc",
			want: []Token{
				{
					Type:  tokenIdent,
					Value: []byte("abc"),
				},
			},
		},
		{
			src: "a b",
			want: []Token{
				{
					Type:  tokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  tokenIdent,
					Value: []byte("b"),
				},
			},
		},
		{
			src: " a b",
			want: []Token{
				{
					Type:  tokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  tokenIdent,
					Value: []byte("b"),
				},
			},
		},
		{
			src: "a    b cd  e fgh",
			want: []Token{
				{
					Type:  tokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  tokenIdent,
					Value: []byte("b"),
				},
				{
					Type:  tokenIdent,
					Value: []byte("cd"),
				},
				{
					Type:  tokenIdent,
					Value: []byte("e"),
				},
				{
					Type:  tokenIdent,
					Value: []byte("fgh"),
				},
			},
		},
		{
			src: "1",
			want: []Token{
				{
					Type:  tokenErr,
					Value: []byte("1"),
				},
			},
		},
		{
			src: `"a"`,
			want: []Token{
				{
					Type:  tokenString,
					Value: []byte(`"a"`),
				},
			},
		},
		{
			src: `a "b"`,
			want: []Token{
				{
					Type:  tokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  tokenString,
					Value: []byte(`"b"`),
				},
			},
		},
		{
			src: `a"b"`,
			want: []Token{
				{
					Type:  tokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  tokenString,
					Value: []byte(`"b"`),
				},
			},
		},
		{
			src: `a "b`,
			want: []Token{
				{
					Type:  tokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  tokenErr,
					Value: []byte(`"b`),
				},
			},
		},
		{
			src: "substitute a b",
			want: []Token{
				{
					Type:  tokenIdent,
					Value: []byte("substitute"),
				},
				{
					Type:  tokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  tokenIdent,
					Value: []byte("b"),
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
