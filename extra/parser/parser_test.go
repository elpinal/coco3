package parser

import (
	"reflect"
	"testing"

	"github.com/elpinal/coco3/extra/ast"
)

func TestParse(t *testing.T) {
	tests := []struct {
		src  string
		name string
		args []ast.Expr
	}{
		{
			src:  "aa 'b'",
			name: "aa",
			args: []ast.Expr{&ast.String{Lit: "b"}},
		},
		{
			src:  "a 'u' : 'v' : []",
			name: "a",
			args: []ast.Expr{
				&ast.Cons{
					Head: "u",
					Tail: &ast.Cons{
						Head: "v",
						Tail: &ast.Empty{},
					},
				},
			},
		},
		{
			src:  "a-b ['u', 'v']",
			name: "a-b",
			args: []ast.Expr{
				&ast.Cons{
					Head: "u",
					Tail: &ast.Cons{
						Head: "v",
						Tail: &ast.Empty{},
					},
				},
			},
		},
		{
			src:  "a-b1-2190 [''] 8",
			name: "a-b1-2190",
			args: []ast.Expr{
				&ast.Cons{
					Head: "",
					Tail: &ast.Empty{},
				},
				&ast.Int{
					Lit: "8",
				},
			},
		},
		{
			src:  `a '{{range .Imports}}{{. | printf "%s\\n"}}{{end}}'`,
			name: "a",
			args: []ast.Expr{
				&ast.String{`{{range .Imports}}{{. | printf "%s\n"}}{{end}}`},
			},
		},
	}
	for _, test := range tests {
		x, err := Parse([]byte(test.src))
		if err != nil {
			t.Fatalf("Parse(%q): %v", test.src, err)
		}
		if x.Name.Lit != test.name {
			t.Fatalf("Parse(%q).Lit != %s; got %s", test.src, test.name, x.Name.Lit)
		}
		if !reflect.DeepEqual(x.Args, test.args) {
			t.Errorf("Parse(%q).Args != %v; got %v", test.src, test.args, x.Args)
		}
	}
}

func TestParseFail(t *testing.T) {
	tests := []string{
		"aa '",
		"'a'",
		"12",
		"a :",
		"a [",
	}
	for _, src := range tests {
		got, err := Parse([]byte(src))
		if err == nil {
			t.Errorf("Parse(%q): unexpectedly succeeded", src)
			t.Fatalf("got: %v", got)
		}
	}
}
