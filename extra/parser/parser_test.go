package parser

import (
	"reflect"
	"testing"

	"github.com/elpinal/coco3/extra/ast"
)

func TestParse(t *testing.T) {
	tests := []struct {
		src  string
		want ast.Command
	}{
		{
			src:  "aa 'b'",
			want: ast.Command{Name: "aa", Args: []ast.Expr{&ast.String{Lit: "b"}}},
		},
		{
			src: "a 'u' : 'v' : []",
			want: ast.Command{Name: "a", Args: []ast.Expr{
				&ast.Cons{
					Head: "u",
					Tail: &ast.Cons{
						Head: "v",
						Tail: &ast.Empty{},
					},
				},
			}},
		},
		{
			src: "a-b ['u', 'v']",
			want: ast.Command{Name: "a-b", Args: []ast.Expr{
				&ast.Cons{
					Head: "u",
					Tail: &ast.Cons{
						Head: "v",
						Tail: &ast.Empty{},
					},
				},
			}},
		},
		{
			src: "a-b1-2190 [''] 8",
			want: ast.Command{Name: "a-b1-2190", Args: []ast.Expr{
				&ast.Cons{
					Head: "",
					Tail: &ast.Empty{},
				},
				&ast.Int{
					Lit: "8",
				},
			}},
		},
	}
	for _, test := range tests {
		x, err := Parse([]byte(test.src))
		if err != nil {
			t.Fatalf("Parse(%q): %v", test.src, err)
		}
		if !reflect.DeepEqual(*x, test.want) {
			t.Errorf("Parse(%q) != %v; got %v", test.src, test.want, x)
		}
	}
}
