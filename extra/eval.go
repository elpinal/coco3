package extra

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/elpinal/coco3/extra/ast"
	"github.com/elpinal/coco3/extra/typed"
	"github.com/elpinal/coco3/extra/types"
)

type Env struct {
	cmds map[string]typed.Command
}

func New() Env {
	return Env{cmds: map[string]typed.Command{
		"exec": execCommand,
		"cd":   cdCommand,
		"exit": exitCommand,
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
	for i, arg := range command.Args {
		if arg.Type() != tc.Params[i] {
			return fmt.Errorf("type mismatch: %v != %v", arg.Type(), tc.Params[i])
		}
	}
	return tc.Fn(command.Args)
}

var execCommand = typed.Command{
	Params: []types.Type{types.String, types.StringList},
	Fn: func(args []ast.Expr) error {
		var cmdArgs []string
	loop:
		for cons := args[1].(*ast.Cons); ; {
			cmdArgs = append(cmdArgs, cons.Head)
			switch x := cons.Tail.(type) {
			case *ast.Empty:
				break loop
			case *ast.Cons:
				cons = x
			}
		}
		return exec.Command(args[0].(*ast.String).Lit, cmdArgs...).Run()
	},
}

var cdCommand = typed.Command{
	Params: []types.Type{types.String},
	Fn: func(args []ast.Expr) error {
		return os.Chdir(args[0].(*ast.String).Lit)
	},
}

var exitCommand = typed.Command{
	Params: []types.Type{types.Int},
	Fn: func(args []ast.Expr) error {
		n, err := strconv.Atoi(args[0].(*ast.Int).Lit)
		if err != nil {
			return err
		}
		os.Exit(n)
		return nil
	},
}
