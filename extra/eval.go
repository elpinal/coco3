package extra

import (
	"fmt"

	"github.com/elpinal/coco3/extra/ast"
	"github.com/elpinal/coco3/extra/typed"
)

type Env struct {
	cmds map[string]typed.Command
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
