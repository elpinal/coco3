package extra

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/elpinal/coco3/extra/ast"
	"github.com/elpinal/coco3/extra/typed"
	"github.com/elpinal/coco3/extra/types"
)

func TestEval(t *testing.T) {
	var buf bytes.Buffer
	prefix := "print: the argument is"
	printCommand := func(args []ast.Expr) error {
		_, err := fmt.Fprintln(&buf, prefix, args[0].(*ast.String).Lit)
		return err
	}
	e := New()
	e.Bind("print", typed.Command{Params: []types.Type{types.String}, Fn: printCommand})
	err := e.Eval(&ast.Command{Name: "print", Args: []ast.Expr{&ast.String{Lit: "aaa"}}})
	if err != nil {
		t.Fatalf("Eval: %v", err)
	}
	got := buf.String()
	want := prefix + " aaa\n"
	if got != want {
		t.Errorf("Eval: want %q, but got %q", want, got)
	}
}
