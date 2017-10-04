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
	if len(tc.params) != 1 {
		return fmt.Errorf("the parameters should be the list of one String: %v", tc.params)
	}
	if tc.params[0] != String {
		return fmt.Errorf("the parameter should be String type: %v", tc.params[0])
	}
	tc.fn(command.Arg.Lit)
	return nil
}
