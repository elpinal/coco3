package parser

import (
	"testing"

	"github.com/elpinal/coco3/extra/ast"
	"github.com/elpinal/coco3/extra/token"
)

func TestParse(t *testing.T) {
	src := "aa 'b'"
	x, err := Parse([]byte(src))
	if err != nil {
		t.Errorf("Parse: %v", err)
	}
	want := ast.Command{Name: "aa", Arg: token.Token{Kind: STRING, Lit: "b"}}
	if *x != want {
		t.Errorf("Parse(%s) != %v; got %v", src, want, x)
	}
}
