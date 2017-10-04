package parser

import (
	"reflect"
	"testing"

	"github.com/elpinal/coco3/extra/ast"
)

func TestParse(t *testing.T) {
	src := "aa 'b'"
	x, err := Parse([]byte(src))
	if err != nil {
		t.Errorf("Parse: %v", err)
	}
	want := ast.Command{Name: "aa", Args: []ast.Expr{&ast.String{Lit: "b"}}}
	if !reflect.DeepEqual(*x, want) {
		t.Errorf("Parse(%s) != %v; got %v", src, want, x)
	}
}
