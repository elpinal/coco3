package extra

import (
	"fmt"
	"os/exec"

	"github.com/elpinal/coco3/extra/ast"
	"github.com/elpinal/coco3/extra/typed"
)

type Env struct {
	cmds map[string]typed.Command
}

func New() Env {
	return Env{cmds: map[string]typed.Command{
		"exec": execCommand,
	}}
}

func WithoutDefault() Env {
	return Env{cmds: make(map[string]typed.Command)}
}

func (e *Env) Bind(name string, c typed.Command) {
	e.cmds[name] = c
}

func (e *Env) Eval(command *ast.Command) error {
	tc, found := e.cmds[command.Name]
	if !found {
		return fmt.Errorf("no such typed command: %q", command.Name)
	}
	if len(command.Args) != len(tc.Params) {
		return fmt.Errorf("the length of args (%d) != the one of params (%d)", len(command.Args), len(tc.Params))
	}
	args := make([]string, 0, len(command.Args))
	for i, arg := range command.Args {
		if arg.Kind != tc.Params[i] {
			return fmt.Errorf("type mismatch: %v != %v", arg.Kind, tc.Params[i])
		}
		args = append(args, arg.Lit)
	}
	return tc.Fn(args)
}

var execCommand = typed.Command{
	Params: []typed.Type{typed.String},
	Fn: func(args []string) error {
		return exec.Command(args[0]).Run()
	},
}
