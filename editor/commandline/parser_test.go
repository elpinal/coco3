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
					Type:  TokenEOF,
					Value: []byte(""),
				},
			},
		},
		{
			src: " ",
			want: []Token{
				{
					Type:  TokenEOF,
					Value: []byte(""),
				},
			},
		},
		{
			src: "a",
			want: []Token{
				{
					Type:  TokenIdent,
					Value: []byte("a"),
				},
			},
		},
		{
			src: "abc",
			want: []Token{
				{
					Type:  TokenIdent,
					Value: []byte("abc"),
				},
			},
		},
		{
			src: "ABc",
			want: []Token{
				{
					Type:  TokenIdent,
					Value: []byte("ABc"),
				},
			},
		},
		{
			src: " abc",
			want: []Token{
				{
					Type:  TokenIdent,
					Value: []byte("abc"),
				},
			},
		},
		{
			src: "a b",
			want: []Token{
				{
					Type:  TokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  TokenIdent,
					Value: []byte("b"),
				},
			},
		},
		{
			src: " a b",
			want: []Token{
				{
					Type:  TokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  TokenIdent,
					Value: []byte("b"),
				},
			},
		},
		{
			src: "a    b cd  e fgh",
			want: []Token{
				{
					Type:  TokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  TokenIdent,
					Value: []byte("b"),
				},
				{
					Type:  TokenIdent,
					Value: []byte("cd"),
				},
				{
					Type:  TokenIdent,
					Value: []byte("e"),
				},
				{
					Type:  TokenIdent,
					Value: []byte("fgh"),
				},
			},
		},
		{
			src: "1",
			want: []Token{
				{
					Type:  TokenErr,
					Value: []byte("1"),
				},
			},
		},
		{
			src: `"a"`,
			want: []Token{
				{
					Type:  TokenString,
					Value: []byte(`"a"`),
				},
			},
		},
		{
			src: `a "b"`,
			want: []Token{
				{
					Type:  TokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  TokenString,
					Value: []byte(`"b"`),
				},
			},
		},
		{
			src: `a"b"`,
			want: []Token{
				{
					Type:  TokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  TokenString,
					Value: []byte(`"b"`),
				},
			},
		},
		{
			src: `a "b`,
			want: []Token{
				{
					Type:  TokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  TokenErr,
					Value: []byte(`"b`),
				},
			},
		},
		{
			src: "substitute a b",
			want: []Token{
				{
					Type:  TokenIdent,
					Value: []byte("substitute"),
				},
				{
					Type:  TokenIdent,
					Value: []byte("a"),
				},
				{
					Type:  TokenIdent,
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

func TestParseT(t *testing.T) {
	got, err := ParseT("substitute a b")
	if err != nil {
		t.Fatalf("ParseT: %v", err)
	}
	want := &CommandT{
		Name: []byte("substitute"),
		Args: []Token{
			{
				Type:  TokenIdent,
				Value: []byte("a"),
			},
			{
				Type:  TokenIdent,
				Value: []byte("b"),
			},
		},
	}
	if !reflect.DeepEqual(got.Name, want.Name) {
		t.Fatalf("ParseT/Name: got %s, want %s", got.Name, want.Name)
	}
	if !reflect.DeepEqual(got.Args, want.Args) {
		t.Fatalf("ParseT/Args: got %v, want %v", got.Args, want.Args)
	}
}
