package parser

import (
	"testing"

	"github.com/elpinal/coco3/extra/token"
)

func TestParse(t *testing.T) {
	src := "aa"
	x, err := Parse([]byte(src))
	if err != nil {
		t.Errorf("Parse: %v", err)
	}
	want := token.Token{Kind: IDENT, Lit: src}
	if x != want {
		t.Errorf("Parse(%s) != %v; got %v", src, want, x)
	}
}
