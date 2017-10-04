package extra

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/elpinal/coco3/extra/parser"
	"github.com/elpinal/coco3/extra/typed"
)

func TestEval(t *testing.T) {
	var buf bytes.Buffer
	prefix := "print: the argument is"
	printCommand := func(args []string) error {
		_, err := fmt.Fprintln(&buf, prefix, args[0])
		return err
	}
	e := Env{cmds: map[string]typed.Command{"print": {Params: []typed.Type{typed.String}, Fn: printCommand}}}
	src := "print 'aaa'"
	c, err := parser.Parse([]byte(src))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	err = e.Eval(c)
	if err != nil {
		t.Fatalf("Eval: %v", err)
	}
	got := buf.String()
	want := prefix + " aaa\n"
	if got != want {
		t.Errorf("Eval(%q) != %q; got %q", src, want, got)
	}
}
