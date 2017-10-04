package extra

import (
	"fmt"

	"github.com/elpinal/coco3/extra/ast"
)

type Env struct {
	cmds map[string]TypedCommand
}

func (e *Env) Eval(command *ast.Command) error {
	tc, found := e.cmds[command.Name]
	if !found {
		return fmt.Errorf("no such typed command: %q", command.Name)
	}
	if len(command.Args) != len(tc.params) {
		return fmt.Errorf("the length of args (%d) != the one of params (%d)", len(command.Args), len(tc.params))
	}
	args := make([]string, 0, len(command.Args))
	for i, arg := range command.Args {
		// FIXME
		if arg.Kind != tc.params[i] {
			return fmt.Errorf("type mismatch: %v != %v", arg.Kind, tc.params[i])
		}
		args = append(args, arg.Lit)
	}
	return tc.fn(args)
}
