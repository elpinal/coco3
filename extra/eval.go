package extra

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/pkg/errors"

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

		"git":   gitCommand,
		"cargo": cargoCommand,
		"go":    goCommand,
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
			return fmt.Errorf("type mismatch: (%v) (type of %v) does not match with (%v) (expected type)", arg.Type(), arg, tc.Params[i])
		}
	}
	return tc.Fn(command.Args)
}

func toSlice(list ast.List) ([]string, error) {
	var ret []string
	for {
		switch x := list.(type) {
		case *ast.Cons:
			ret = append(ret, x.Head)
			list = x.Tail
		case *ast.Empty:
			return ret, nil
		default:
			return nil, fmt.Errorf("unexpected list type: %T", x)
		}
	}
}

var execCommand = typed.Command{
	Params: []types.Type{types.String, types.StringList},
	Fn: func(args []ast.Expr) error {
		cmdArgs, err := toSlice(args[1].(ast.List))
		if err != nil {
			return errors.Wrap(err, "exec")
		}
		cmd := exec.Command(args[0].(*ast.String).Lit, cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
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

func commandsInCommand(name string) func([]ast.Expr) error {
	return func(args []ast.Expr) error {
		cmdArgs, err := toSlice(args[1].(ast.List))
		if err != nil {
			return errors.Wrap(err, name)
		}
		var cmd *exec.Cmd
		switch lit := args[0].(*ast.Ident).Lit; lit {
		case "command":
			cmd = exec.Command(name, cmdArgs...)
		default:
			cmd = exec.Command(name, append([]string{lit}, cmdArgs...)...)
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}
}

var gitCommand = typed.Command{
	Params: []types.Type{types.Ident, types.StringList},
	Fn:     commandsInCommand("git"),
}

var cargoCommand = typed.Command{
	Params: []types.Type{types.Ident, types.StringList},
	Fn:     commandsInCommand("cargo"),
}

var goCommand = typed.Command{
	Params: []types.Type{types.Ident, types.StringList},
	Fn:     commandsInCommand("go"),
}
