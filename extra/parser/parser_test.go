package parser

import (
	"reflect"
	"testing"

	"github.com/elpinal/coco3/extra/ast"
	"github.com/elpinal/coco3/extra/token"
	"github.com/elpinal/coco3/extra/typed"
)

func TestParse(t *testing.T) {
	src := "aa 'b'"
	x, err := Parse([]byte(src))
	if err != nil {
		t.Errorf("Parse: %v", err)
	}
	want := ast.Command{Name: "aa", Args: []token.Token{{Kind: typed.String, Lit: "b"}}}
	if !reflect.DeepEqual(*x, want) {
		t.Errorf("Parse(%s) != %v; got %v", src, want, x)
	}
}
